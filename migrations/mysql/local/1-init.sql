DROP DATABASE IF EXISTS db_bemax_api;

CREATE DATABASE IF NOT EXISTS db_bemax_api;
USE db_bemax_api;

-- TABELAS DE USUÁRIOS E AUTENTICAÇÃO

DROP TABLE IF EXISTS USER_ROLES;
DROP TABLE IF EXISTS ADDRESSES;
DROP TABLE IF EXISTS LOGIN_ATTEMPTS;
DROP TABLE IF EXISTS TOKENS;
DROP TABLE IF EXISTS ACTIVE_SESSIONS;
DROP TABLE IF EXISTS OAUTH_ACCOUNTS;
DROP TABLE IF EXISTS ROLES;
DROP TABLE IF EXISTS USERS;
DROP TABLE IF EXISTS STATES;

CREATE TABLE USERS (
   id CHAR(36) NOT NULL PRIMARY KEY,
   email VARCHAR(255) NOT NULL UNIQUE,
   password_hash VARCHAR(255) NULL,
   auth_provider ENUM('local', 'oauth') NOT NULL DEFAULT 'local',
   full_name VARCHAR(180) NOT NULL,
   cpf VARCHAR(14) NULL UNIQUE,
   phone VARCHAR(20) NULL,
   birth_date DATE NULL,
   profile_picture VARCHAR(500) NULL,
   email_verified BOOLEAN NOT NULL DEFAULT FALSE,
   phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
   profile_completed BOOLEAN NOT NULL DEFAULT FALSE,
   last_login TIMESTAMP NULL,
   status ENUM('active', 'inactive', 'blocked', 'pending_verification') NOT NULL DEFAULT 'active',
   token_version INT NOT NULL DEFAULT 0,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   INDEX idx_users_email (email),
   INDEX idx_users_cpf (cpf),
   INDEX idx_users_auth_provider (auth_provider),
   INDEX idx_users_status (status),
   INDEX idx_users_profile_completed (profile_completed),
   INDEX idx_users_token_version (id, token_version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE OAUTH_ACCOUNTS (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    provider VARCHAR(20) NOT NULL,
    provider_uid VARCHAR(255) NOT NULL,
    firebase_uid VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    provider_picture VARCHAR(500) NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
    UNIQUE KEY unique_provider_uid (provider, provider_uid),
    UNIQUE KEY unique_firebase_uid (firebase_uid),
    INDEX idx_oauth_user_id (user_id),
    INDEX idx_oauth_provider (provider),
    INDEX idx_provider_uid (provider, provider_uid),
    INDEX idx_oauth_provider_email (provider_email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE ACTIVE_SESSIONS (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    session_id VARCHAR(100) NOT NULL UNIQUE,
    last_access_token_jti VARCHAR(255) NOT NULL,
    device_type ENUM('mobile', 'desktop', 'tablet', 'unknown') NULL,
    user_agent TEXT NULL,
    ip_address VARCHAR(50) NULL,
    is_suspicious BOOLEAN NOT NULL DEFAULT FALSE,
    risk_score INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_refreshed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
    INDEX idx_sessions_user (user_id),
    INDEX idx_sessions_session_id (session_id),
    INDEX idx_sessions_jti (last_access_token_jti),
    INDEX idx_sessions_active (user_id, is_active, expires_at),
    INDEX idx_sessions_suspicious (is_suspicious, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE ROLES (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_role_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE USER_ROLES (
    user_id CHAR(36) NOT NULL,
    role_id CHAR(36) NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES ROLES(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE STATES (
    id CHAR(2) PRIMARY KEY UNIQUE,
    name VARCHAR(50) NOT NULL UNIQUE,
region ENUM('Norte', 'Nordeste', 'Centro-Oeste', 'Sudeste', 'Sul') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE ADDRESSES (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    street VARCHAR(255) NOT NULL,
    number VARCHAR(20) NOT NULL,
    complement VARCHAR(255) NULL,
    neighborhood VARCHAR(100) NULL,
    city VARCHAR(100) NOT NULL,
    state_id CHAR(2) NOT NULL,
    zip_code VARCHAR(10) NOT NULL,
    latitude DECIMAL(10, 8) NULL,
    longitude DECIMAL(11, 8) NULL,
    is_default BOOLEAN DEFAULT FALSE,
    type ENUM('residential', 'commercial', 'billing', 'shipping') NOT NULL DEFAULT 'residential',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
    FOREIGN KEY (state_id) REFERENCES STATES(id),
    INDEX idx_addresses_user (user_id),
    INDEX idx_addresses_zip_code (zip_code),
    INDEX idx_addresses_default (user_id, is_default),
    INDEX idx_addresses_location (latitude, longitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE TOKENS (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    session_id CHAR(36) NULL, -- TODO Analyze whether it needs to be NOT NULL.
    token VARCHAR(255) NOT NULL,
    token_type VARCHAR(50) NOT NULL, -- TODO Validar se precisa salvar mesmo, sem utilidade no momento...
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revoked_at TIMESTAMP NULL,
    revoked_reason VARCHAR(100) NULL COMMENT 'logout, security, expired, user_request',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES ACTIVE_SESSI ONS(id) ON DELETE SET NULL,
    INDEX idx_tokens_user (user_id),
    INDEX idx_tokens_session (session_id),
    INDEX idx_tokens_token (token),
    INDEX idx_tokens_revoked (is_revoked, expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Table for login attempts
CREATE TABLE LOGIN_ATTEMPTS (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    success BOOLEAN NOT NULL,
    ip_address VARCHAR(50) NULL,
    user_agent TEXT NULL,
    failure_reason VARCHAR(100) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_LOGIN_ATTEMPTS_EMAIL (email),
    INDEX idx_LOGIN_ATTEMPTS_CREATED_AT (created_at),
    INDEX idx_LOGIN_ATTEMPTS_IP (ip_address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;