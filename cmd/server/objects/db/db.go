package db

import (
	"database/sql"
	"errors"
	"gosecureskeleton/cmd/server/objects"
	"os"
	"strings"
	"time"
)

type DBStore struct {
	db *sql.DB
}

func (s *DBStore) Close() error {
	return s.db.Close()
}

func (s *DBStore) OpenStore(databasePath, schemaFile string) error {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(1)

	s.db = db
	if err := s.Initialize(schemaFile); err != nil {
		_ = db.Close()
		return err
	}

	return nil
}

func (s *DBStore) Initialize(schemaFile string) error {
	if err := s.ExecSQLFile(schemaFile); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) ExecSQLFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(string(content))
	return err
}

func (s *DBStore) FindUserByUsername(username string) (objects.User, bool, error) {
	row := s.db.QueryRow(`
		SELECT id, username, name, email, phone, password, password_salt, balance, is_admin
		FROM users
		WHERE username = ?
	`, strings.TrimSpace(username))

	var user objects.User
	var isAdmin int64
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Salt, &user.Balance, &isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return objects.User{}, false, nil
		}
		return objects.User{}, false, err
	}
	user.IsAdmin = isAdmin == 1

	return user, true, nil
}

func (s *DBStore) InsertUser(user objects.User) error {
	const query string = "INSERT INTO users (username, name, email, phone, password, password_salt, balance, is_admin) VALUES (?, ?, ?, ?, ?, ?, 0, ?);"
	if _, err := s.db.Exec(query, strings.TrimSpace(user.Username), user.Name, user.Email, user.Phone, user.Password, user.Salt, user.Balance, user.IsAdmin); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) DeleteUserByID(ID uint) error {
	const query string = "DELETE FROM users WHERE id = ?;"
	if _, err := s.db.Exec(query, ID); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) UpdateUser(user objects.User) error {
	const query string = "UPDATE users SET name = ?, email = ?, phone = ?, is_admin = ? WHERE id = ?;"
	if _, err := s.db.Exec(query, user.Name, user.Email, user.Phone, user.IsAdmin, user.ID); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) AddUserBalenceByID(ID uint, amount int64) (bool, error) {
	const userBalenceCheckQuery string = "SELECT balance FROM users WHERE id = ?;"
	row := s.db.QueryRow(userBalenceCheckQuery, ID)
	var ABalence int64
	if err := row.Scan(&ABalence); err != nil {
		return false, err
	}
	if ABalence+amount < 0 {
		return false, nil
	}

	const query string = "UPDATE users SET balance = balance + ? WHERE id = ?;"
	if _, err := s.db.Exec(query, amount, ID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *DBStore) TransferBalenceAToB(AID uint, BID uint, amount int64) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	const userBalenceCheckQuery string = "SELECT balance FROM users WHERE id = ?;"
	row := tx.QueryRow(userBalenceCheckQuery, AID)
	var ABalence int64
	if err := row.Scan(&ABalence); err != nil {
		return false, err
	}
	if ABalence-amount < 0 {
		return false, nil
	}

	const updateMinusAccountQuery string = "UPDATE users SET balance = balance - ? WHERE id = ?"
	_, err = tx.Exec(updateMinusAccountQuery, amount, AID)
	if err != nil {
		return false, err
	}
	const updatePlusAccountQuery string = "UPDATE users SET balance = balance + ? WHERE id = ?"
	_, err = tx.Exec(updatePlusAccountQuery, amount, BID)
	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (s *DBStore) InsertPage(post objects.PostView) error {
	const query string = "INSERT INTO posts (title, content, owner_id) VALUES (?, ?, ?);"
	if _, err := s.db.Exec(query, post.Title, post.Content, post.OwnerID); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) UpdatePage(post objects.PostView) error {
	const query string = "UPDATE posts SET title = ?, content = ?, owner_id = ?, updated_at = ? WHERE id = ?;"
	if _, err := s.db.Exec(query, post.Title, post.Content, post.OwnerID, time.Now(), post.ID); err != nil {
		return err
	}
	return nil
}

var DB DBStore
