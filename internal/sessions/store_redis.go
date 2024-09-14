package sessions

import (
	"context"
	"crypto/tls"
	"encoding/base32"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/boj/redistore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	redis "github.com/redis/go-redis/v9"
)

// ATTRIBUTION: This implementation takes heavy inspiration from
// https://github.com/boj/redistore, while replacing the underlying Redigo library
// with go-redis in order to support secure connections.

// Amount of time for cookies/redis keys to expire.
var sessionExpire = 86400 * 30

const (
	keyPrefix     = "session:"
	maxLength     = 4096
	defaultMaxAge = 60 * 20 // 20 minutes seems like a reasonable default
)

type store struct {
	client         *redis.Client
	Codecs         []securecookie.Codec
	serializer     redistore.SessionSerializer
	SessionOptions *sessions.Options
}

// NewRedis instantiates a new client
func NewRedis(addr string, keyPairs ...[]byte) (sessions.Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client := redis.NewClient(&redis.Options{Addr: addr})
	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", status.Err())
	}

	return &store{
		client:     client,
		serializer: redistore.GobSerializer{},
		Codecs:     securecookie.CodecsFromPairs(keyPairs...),
		SessionOptions: &sessions.Options{
			Path:   "/",
			MaxAge: sessionExpire,
		},
	}, nil
}

// NewRedisSecure instantiates a new secure TLS client
func NewRedisSecure(host, port, user, password string, db int, keyPairs ...[]byte) (sessions.Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", host, port),
		Username:  user,
		Password:  password,
		TLSConfig: &tls.Config{},
		DB:        db,
	})
	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, fmt.Errorf("failed to create secure Redis client: %w", status.Err())
	}

	return &store{
		client:     client,
		serializer: redistore.GobSerializer{},
		Codecs:     securecookie.CodecsFromPairs(keyPairs...),
		SessionOptions: &sessions.Options{
			Path:   "/",
			MaxAge: sessionExpire,
		},
	}, nil
}

// Get should return a cached session.
func (s *store) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New should create and return a new session.
//
// Note that New should never return a nil session, even in the case of
// an error if using the Registry infrastructure to cache the session.
func (s *store) New(r *http.Request, name string) (*sessions.Session, error) {
	var (
		err error
		ok  bool
	)
	session := sessions.NewSession(s, name)
	// make a copy
	options := *s.SessionOptions
	session.Options = &options
	session.IsNew = true
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {
			ok, err = s.load(session)
			session.IsNew = !(err == nil && ok) // not new if no error and data available
		}
	}
	return session, err
}

// Save should persist session to the underlying store implementation.
func (s *store) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Marked for deletion.
	if session.Options.MaxAge <= 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		if err := s.save(session); err != nil {
			return err
		}
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
		if err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	}
	return nil
}

func (s *store) Options(options sessions.Options) {
	s.SessionOptions = &options
}

// save stores the session in redis.
func (s *store) save(session *sessions.Session) error {
	b, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}
	if maxLength != 0 && len(b) > maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	age := session.Options.MaxAge
	if age == 0 {
		age = defaultMaxAge
	}
	cmd := s.client.SetEx(context.TODO(), keyPrefix+session.ID, b, time.Duration(age*int(time.Second)))
	return cmd.Err()
}

// load reads the session from redis.
// returns true if there is session data in DB
func (s *store) load(session *sessions.Session) (bool, error) {
	cmd := s.client.Get(context.TODO(), keyPrefix+session.ID)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}

	if cmd.String() == "" {
		return false, nil // no data was associated with this key
	}

	b, err := cmd.Bytes()
	if err != nil {
		return false, err
	}
	return true, s.serializer.Deserialize(b, session)
}

// delete removes keys from redis
func (s *store) delete(session *sessions.Session) error {
	cmd := s.client.Del(context.TODO(), keyPrefix+session.ID)

	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}
