CREATE DATABASE IF NOT EXISTS feishu_suite CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE feishu_suite;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY,
    union_id VARCHAR(64) UNIQUE,
    open_id VARCHAR(64) UNIQUE,
    name VARCHAR(128) NOT NULL,
    en_name VARCHAR(128),
    email VARCHAR(128),
    phone VARCHAR(32),
    avatar_url VARCHAR(512),
    avatar_thumb VARCHAR(512),
    avatar_middle VARCHAR(512),
    status VARCHAR(32) DEFAULT 'active',
    is_activated BOOLEAN DEFAULT TRUE,
    is_tenant_access BOOLEAN DEFAULT TRUE,
    department_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_phone (phone),
    INDEX idx_union_id (union_id),
    INDEX idx_department_id (department_id)
);

CREATE TABLE IF NOT EXISTS departments (
    department_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    name_en VARCHAR(128),
    parent_id VARCHAR(64),
    `order` INT DEFAULT 0,
    is_root BOOLEAN DEFAULT FALSE,
    member_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id VARCHAR(64) PRIMARY KEY,
    chat_id VARCHAR(64) NOT NULL,
    sender VARCHAR(64),
    sender_id VARCHAR(64),
    sender_type VARCHAR(32),
    content TEXT,
    msg_type VARCHAR(32),
    is_deleted BOOLEAN DEFAULT FALSE,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_chat_id (chat_id),
    INDEX idx_sender_id (sender_id)
);

CREATE TABLE IF NOT EXISTS chats (
    chat_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    member_count INT DEFAULT 0,
    owner_id VARCHAR(64),
    owner_id_type VARCHAR(32),
    is_activated BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_owner_id (owner_id)
);

CREATE TABLE IF NOT EXISTS calendars (
    calendar_id VARCHAR(64) PRIMARY KEY,
    summary VARCHAR(128) NOT NULL,
    description TEXT,
    type VARCHAR(32),
    timezone VARCHAR(64),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS events (
    event_id VARCHAR(64) PRIMARY KEY,
    calendar_id VARCHAR(64) NOT NULL,
    summary VARCHAR(256) NOT NULL,
    description TEXT,
    is_all_day BOOLEAN DEFAULT FALSE,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    timezone VARCHAR(64),
    attendees JSON,
    status VARCHAR(32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_calendar_id (calendar_id),
    INDEX idx_start_time (start_time)
);

CREATE TABLE IF NOT EXISTS approvals (
    approval_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    form_content JSON,
    instance_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS approval_instances (
    instance_id VARCHAR(64) PRIMARY KEY,
    approval_id VARCHAR(64) NOT NULL,
    title VARCHAR(256) NOT NULL,
    status VARCHAR(32),
    initiator_id VARCHAR(64),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_approval_id (approval_id),
    INDEX idx_status (status)
);

CREATE TABLE IF NOT EXISTS casbin_rule (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    ptype VARCHAR(128) NOT NULL,
    v0 VARCHAR(128) NOT NULL,
    v1 VARCHAR(128) NOT NULL,
    v2 VARCHAR(128) NOT NULL,
    v3 VARCHAR(128),
    v4 VARCHAR(128),
    v5 VARCHAR(128),
    INDEX idx_ptype (ptype),
    INDEX idx_v0 (v0),
    INDEX idx_v1 (v1)
);

INSERT INTO users (id, union_id, open_id, name, email, status, is_activated) VALUES
('admin', 'admin_union_id', 'admin_open_id', 'System Admin', 'admin@example.com', 'active', TRUE);

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'admin', '/api/v1/*', 'GET'),
('p', 'admin', '/api/v1/*', 'POST'),
('p', 'admin', '/api/v1/*', 'PUT'),
('p', 'admin', '/api/v1/*', 'DELETE'),
('p', 'user', '/api/v1/user', 'GET'),
('p', 'user', '/api/v1/message', 'GET'),
('p', 'user', '/api/v1/message', 'POST');