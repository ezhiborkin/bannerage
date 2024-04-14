CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   email VARCHAR(30) UNIQUE NOT NULL,
   role VARCHAR(10) NOT NULL DEFAULT 'user',
   password_hash VARCHAR(60) NOT NULL
);

CREATE TABLE banners (
    banner_id SERIAL PRIMARY KEY,
    chosen_revision_id INT DEFAULT NULL
);

CREATE TABLE banner_revisions (
     revision_id SERIAL PRIMARY KEY,
     banner_id INT NOT NULL,
     feature_id INT NOT NULL,
     is_active BOOL DEFAULT TRUE,
     content JSONB,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание внешнего ключа для banner_id
ALTER TABLE banner_revisions
ADD CONSTRAINT fk_banner_id
FOREIGN KEY (banner_id)
REFERENCES banners(banner_id)
ON DELETE CASCADE;

-- Создание внешнего ключа для revision_id
ALTER TABLE banners
ADD CONSTRAINT fk_chosen_revision_revision_id
FOREIGN KEY (chosen_revision_id)
REFERENCES banner_revisions(revision_id)
ON DELETE SET NULL;

CREATE TABLE revision_tags (
   revision_id INT NOT NULL REFERENCES banner_revisions(revision_id) ON DELETE CASCADE,
   tag_id INT NOT NULL,
   PRIMARY KEY (revision_id, tag_id)
);