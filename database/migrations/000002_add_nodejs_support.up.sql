-- master table
CREATE TABLE
    IF NOT EXISTS stacks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid VARCHAR(255) NOT NULL UNIQUE,
        name VARCHAR(150) NOT NULL,
        slug VARCHAR(150) NOT NULL,
        type VARCHAR(20) NOT NULL,
        repo_url TEXT,
        branch VARCHAR(150) DEFAULT 'master',
        remote VARCHAR(100) DEFAULT 'origin',
        port INTEGER NOT NULL CHECK (
            port >= 1024
            AND port <= 65535
        ), -- 1024–65535
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

-- pm2-specific configuration for nodejs apps
CREATE TABLE
    IF NOT EXISTS pm2_configs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stack_id INTEGER NOT NULL,
        script VARCHAR(255) NOT NULL,
        name VARCHAR(150) NOT NULL,
        args TEXT,
        watch BOOLEAN DEFAULT 0,
        instances INTEGER DEFAULT 1,
        FOREIGN KEY (stack_id) REFERENCES stacks (id) ON DELETE CASCADE
    );

-- Nginx configuration
CREATE TABLE
    IF NOT EXISTS nginx_configs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stack_id INTEGER NOT NULL,
        domain VARCHAR(255) NOT NULL,
        port INTEGER NOT NULL,
        ssl_enabled BOOLEAN DEFAULT 0,
        custom_conf TEXT,
        FOREIGN KEY (stack_id) REFERENCES stacks (id) ON DELETE CASCADE
    );

-- Deployment history with rollback support
CREATE TABLE
    IF NOT EXISTS deployments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stack_id INTEGER NOT NULL,
        status VARCHAR(20) NOT NULL,
        commit_hash VARCHAR(255),
        log TEXT,
        is_rollback BOOLEAN DEFAULT 0,
        rolled_back_from_id INTEGER,
        deployed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (stack_id) REFERENCES stacks (id) ON DELETE CASCADE,
        FOREIGN KEY (rolled_back_from_id) REFERENCES deployments (id) ON DELETE SET NULL
    );

-- Environment variables
CREATE TABLE
    IF NOT EXISTS env_vars (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stack_id INTEGER NOT NULL,
        key VARCHAR(255) NOT NULL,
        value TEXT NOT NULL,
        is_secret BOOLEAN DEFAULT 0,
        FOREIGN KEY (stack_id) REFERENCES stacks (id) ON DELETE CASCADE
    );
