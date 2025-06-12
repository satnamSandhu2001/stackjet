-- init.sql
-- users
CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email VARCHAR(255) NOT NULL UNIQUE,
        password VARCHAR(255) NOT NULL,
        role VARCHAR(20) NOT NULL
    );

-- master table
CREATE TABLE
    IF NOT EXISTS stacks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid VARCHAR(255) NOT NULL UNIQUE,
        name VARCHAR(150) NOT NULL,
        directory TEXT NOT NULL UNIQUE,
        type VARCHAR(20) NOT NULL,
        repo_url TEXT NOT NULL,
        branch VARCHAR(150) DEFAULT 'master',
        remote VARCHAR(100) DEFAULT 'origin',
        port INTEGER NOT NULL,
        commands TEXT NOT NULL,
        created_successfully BOOLEAN DEFAULT 0,
        initial_deployment_success BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

-- pm2-specific configuration for nodejs apps
CREATE TABLE
    IF NOT EXISTS pm2_configs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        stack_id INTEGER NOT NULL,
        script VARCHAR(255) NOT NULL,
        name VARCHAR(150) NOT NULL,
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
        rolled_back_from_id INTEGER DEFAULT NULL,
        deployed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (stack_id) REFERENCES stacks (id) ON DELETE CASCADE,
        FOREIGN KEY (rolled_back_from_id) REFERENCES deployments (id) ON DELETE SET NULL
    );

-- Deployment logs
CREATE TABLE
    IF NOT EXISTS deployment_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        deployment_id INTEGER NOT NULL,
        log TEXT NOT NULL,
        FOREIGN KEY (deployment_id) REFERENCES deployments (id) ON DELETE CASCADE
    );
