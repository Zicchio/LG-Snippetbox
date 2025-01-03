package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zicchio/LG-Snippetbox/pkg/models"
	"github.com/Zicchio/LG-Snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql" // imported as we need the init() function to register mysql in database/sql
	"github.com/golangcollege/sessions"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

const (
	TlsCertificatePath = "./tls/cert.pem"
	TlsSecretKeyPath   = "./tls/key.pem"
)

type snippetModeler interface {
	Insert(string, string, string) (int, error)
	Get(int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type userModeler interface {
	Insert(string, string, string) error
	Authenticate(string, string) (int, error)
	Get(int) (*models.User, error)
}

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      snippetModeler
	templateCache map[string]*template.Template
	users         userModeler
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // db.Ping() is used to verify that a connection was established as connections are lazy
		return nil, err
	}
	return db, nil
}

// ping is a demo function used to introduce testing
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	addr := flag.String("addr", ":4000", "Http Network Address")
	dns := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name") // NOTE: in non-demo environment, use a real password instead of pass
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")          // NOTE: having a default secret key in the code is clearly bad security practice and is done only for demonstration purposes
	flag.Parse()

	// add a logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dns)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	// Prevention measure against CSRF attacks
	// session.Secure = true
	// session.SameSite = http.SameSiteStrictMode

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}, // NOTE: only modern encryption are allowed: old browsers might not comply
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12, // NOTE: prevent TLS downgrade attack, but old software might not comply
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,  // NOTE: slow read timeouts prevents slow-clients attacks (such as Slowloris)
		WriteTimeout: 10 * time.Second, // NOTE: behaviour of WriteTimeout changes between HTTP and HTTPS servers
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(TlsCertificatePath, TlsSecretKeyPath)
	errorLog.Fatal(err)
}
