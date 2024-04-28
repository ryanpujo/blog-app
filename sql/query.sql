-- Users table
CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('John', 'Doe', 'johndoe', 'password123', 'john.doe@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Jane', 'Smith', 'janesmith', 'password123', 'jane.smith@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Emily', 'Johnson', 'emilyj', 'password123', 'emily.johnson@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Michael', 'Brown', 'michaelb', 'password123', 'michael.brown@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('David', 'Williams', 'davidw', 'password123', 'david.williams@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Sarah', 'Miller', 'sarahm', 'password123', 'sarah.miller@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Daniel', 'Davis', 'danield', 'password123', 'daniel.davis@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Jessica', 'Garcia', 'jessicag', 'password123', 'jessica.garcia@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('William', 'Martinez', 'williamm', 'password123', 'william.martinez@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Laura', 'Hernandez', 'laurah', 'password123', 'laura.hernandez@example.com');



-- Roles table
CREATE TABLE public.roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

-- Permissions table
CREATE TABLE public.permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

-- User_roles table to associate users with roles
CREATE TABLE public.user_roles (
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES public.users(id),
    FOREIGN KEY (role_id) REFERENCES public.roles(id)
);

-- Role_permissions table to associate roles with permissions
CREATE TABLE public.role_permissions (
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES public.roles(id),
    FOREIGN KEY (permission_id) REFERENCES public.permissions(id)
);

-- stories table
CREATE TABLE public.stories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL REFERENCES users(id),
    slug VARCHAR(255) UNIQUE,
    excerpt TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMP WITH TIME ZONE,
    type ENUM('flash_fiction', 'short_story', 'novelette', 'novella') NOT NULL, 
    word_count INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE public.stories_categories (
    story_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (story_id, category_id),
    FOREIGN KEY (story_id) REFERENCES public.stories(id),
    FOREIGN KEY (category_id) REFERENCES public.categories(id)
);

-- User Follow System
CREATE TABLE public.user_follows (
    follower_id INT NOT NULL,
    followed_id INT NOT NULL,
    PRIMARY KEY (follower_id, followed_id),
    FOREIGN KEY (follower_id) REFERENCES public.users(id),
    FOREIGN KEY (followed_id) REFERENCES public.users(id)
);

-- Image Storage
CREATE TABLE public.images (
    id SERIAL PRIMARY KEY,
    story_id INT NOT NULL,
    image_url TEXT NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (story_id) REFERENCES public.stories(id)
);

-- Likes table
CREATE TABLE public.likes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    story_id INT NOT NULL REFERENCES stories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Comments table (simplified for hierarchical modeling)
CREATE TABLE public.comments (
    id SERIAL PRIMARY KEY,
    story_id INT NOT NULL REFERENCES stories(id),
    user_id INT NOT NULL REFERENCES users(id),
    parent_comment_id INT REFERENCES comments(id), -- Self-referencing 
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tags table
CREATE TABLE public.tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Post_tags table to associate posts with tags (many-to-many)
CREATE TABLE public.post_tags (
    story_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (story_id, tag_id),
    FOREIGN KEY (story_id) REFERENCES public.stories(id),
    FOREIGN KEY (tag_id) REFERENCES public.tags(id)
);


-- Indexes for optimization
CREATE INDEX idx_users_username ON public.users(username);
CREATE INDEX idx_stories_author_id ON public.stories(author_id);
CREATE INDEX idx_images_story_id ON public.images(story_id);
CREATE INDEX idx_stories_slug ON public.stories(slug);
CREATE INDEX idx_stories_published_at ON public.stories(published_at);

-- Triggers for automatic 'updated_at' timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_modtime
BEFORE UPDATE ON public.users
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Trigger for stories table
CREATE TRIGGER update_story_modtime
BEFORE UPDATE ON public.stories
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Trigger for roles table
CREATE TRIGGER update_role_modtime
BEFORE UPDATE ON public.roles
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Trigger for permissions table
CREATE TRIGGER update_permission_modtime
BEFORE UPDATE ON public.permissions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Trigger for user_roles table
CREATE TRIGGER update_user_role_modtime
BEFORE UPDATE ON public.user_roles
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- Trigger for role_permissions table
CREATE TRIGGER update_role_permission_modtime
BEFORE UPDATE ON public.role_permissions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
