package ponzu_driver_bolt_db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/fanky5g/ponzu-driver-bolt-db/analytics"
	"github.com/fanky5g/ponzu-driver-bolt-db/config"
	"github.com/fanky5g/ponzu-driver-bolt-db/content"
	"github.com/fanky5g/ponzu-driver-bolt-db/credential"
	recoverykey "github.com/fanky5g/ponzu-driver-bolt-db/recovery-key"
	"github.com/fanky5g/ponzu-driver-bolt-db/uploads"
	"github.com/fanky5g/ponzu-driver-bolt-db/users"
	ponzuConfig "github.com/fanky5g/ponzu/config"
	ponzuContent "github.com/fanky5g/ponzu/content"
	ponzuDriver "github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"path/filepath"
)

type driver struct {
	store        *bolt.DB
	repositories map[tokens.Repository]interface{}
}

func (d *driver) Get(token tokens.Repository) interface{} {
	if repository, ok := d.repositories[token]; ok {
		return repository
	}

	log.Fatalf("Repository %s is not implemented", token)
	return nil
}

func (d *driver) Close() error {
	return d.store.Close()
}

func New(contentTypes map[string]ponzuContent.Builder) (ponzuDriver.Database, error) {
	cfg, err := ponzuConfig.New()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}

	store, err := bolt.Open(filepath.Join(cfg.Paths.DataDir, "system.db"), 0666, nil)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	configRepository, err := config.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize config repository %v", err)
	}

	analyticsRepository, err := analytics.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize analytics repository: %v", err)
	}

	userRepository, err := users.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}

	contentRepository, err := content.New(store, contentTypes)
	if err != nil {
		log.Fatalf("Failed to initialize entities repository: %v", err)
	}

	credentialRepository, err := credential.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize credential repository: %v", err)
	}

	recoveryKeyRepository, err := recoverykey.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize recovery key repository: %v", err)
	}

	uploadRepository, err := uploads.New(store)
	if err != nil {
		log.Fatalf("Failed to initialize upload repository: %v", err)
	}
	// End initialize repositories

	repos := make(map[tokens.Repository]interface{})
	repos[tokens.AnalyticsRepositoryToken] = analyticsRepository
	repos[tokens.ConfigRepositoryToken] = configRepository
	repos[tokens.UserRepositoryToken] = userRepository
	repos[tokens.ContentRepositoryToken] = contentRepository
	repos[tokens.RecoveryKeyRepositoryToken] = recoveryKeyRepository
	repos[tokens.UploadRepositoryToken] = uploadRepository
	repos[tokens.CredentialHashRepositoryToken] = credentialRepository

	return &driver{store: store, repositories: repos}, nil
}
