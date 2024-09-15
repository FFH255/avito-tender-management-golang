DROP TABLE IF EXISTS bid_snapshot;
DROP TABLE IF EXISTS bid;
DROP TABLE IF EXISTS decision;
DROP TABLE IF EXISTS tender_snapshot;
DROP TABLE IF EXISTS tender;
DROP TABLE IF EXISTS organization_responsible;
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS organization;

CREATE TABLE IF NOT EXISTS employee (
    id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id VARCHAR(100) PRIMARY KEY,
    organization_id VARCHAR(100) REFERENCES organization(id) ON DELETE CASCADE,
    user_id VARCHAR(100) REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tender (
    id VARCHAR(100),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type VARCHAR(100) NOT NULL,
    status VARCHAR(100) NOT NULL,
    organization_id VARCHAR(100),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tender_snapshot (
    id VARCHAR(100),
    tender_id VARCHAR(100),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type VARCHAR(100) NOT NULL,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid (
    id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status VARCHAR(100) NOT NULL,
    tender_id VARCHAR(100) NOT NULL,
    author_type VARCHAR(100) NOT NULL,
    author_id VARCHAR(100) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid_snapshot (
    id VARCHAR(100) NOT NULL,
    bid_id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS decision (
    id VARCHAR(100) NOT NULL,
    author_id VARCHAR(100) NOT NULL,
    bid_id VARCHAR(100) NOT NULL,
    tender_id VARCHAR(100) NOT NULL,
    status VARCHAR(100)
);

-- Insert mock data into employee table
INSERT INTO employee (id, username, first_name, last_name)
VALUES
    ('b93e6f53-8c2a-4e32-8a9d-c3498e85cb81', 'user1', 'user1', 'user1'),
    ('a97e7e41-6c7d-4c6b-9894-019f9f2f021c', 'user2', 'user2', 'user2'),
    ('f5f8bb13-d799-4d12-8b45-814d3b8ad3e1', 'user3', 'user3', 'user3');

-- Insert mock data into organization table
INSERT INTO organization (id, name, description, type)
VALUES
    ('719c642e-4900-4b45-bdfd-b4f7d1d2d093', 'Tech Corp', 'A technology company', 'IE'),
    ('7abdf89a-8eb3-438b-8391-64f682ec3b94', 'Health Inc', 'A healthcare provider', 'LLC'),
    ('4d6e5f53-5f14-4a0a-b230-6f0d5c8c2a5b', 'EduNation', 'An educational organization', 'JSC');

-- Insert mock data into organization_responsible table
INSERT INTO organization_responsible (id, organization_id, user_id)
VALUES
    ('45f5e6d3-84eb-4d9f-9f1b-870245cc1234', '719c642e-4900-4b45-bdfd-b4f7d1d2d093', 'b93e6f53-8c2a-4e32-8a9d-c3498e85cb81'),
    ('5e7d6b97-b14a-49a6-b849-69c8b91a1822', '7abdf89a-8eb3-438b-8391-64f682ec3b94', 'a97e7e41-6c7d-4c6b-9894-019f9f2f021c'),
    ('6e7a5b4d-72b3-49aa-848a-17e05a4f56ea', '4d6e5f53-5f14-4a0a-b230-6f0d5c8c2a5b', 'f5f8bb13-d799-4d12-8b45-814d3b8ad3e1');
