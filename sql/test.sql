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

-- Blogs table
CREATE TABLE public.blogs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL REFERENCES users(id),
    slug VARCHAR(255) UNIQUE,
    excerpt TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- Insert queries for 'users' table
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('John', 'Doe', 'johndoe', 'password123', 'john.doe@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Jane', 'Smith', 'janesmith', 'password123', 'jane.smith@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Emily', 'Johnson', 'emilyjohnson', 'password123', 'emily.johnson@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Michael', 'Brown', 'michaelbrown', 'password123', 'michael.brown@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('David', 'Jones', 'davidjones', 'password123', 'david.jones@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Sarah', 'Miller', 'sarahmiller', 'password123', 'sarah.miller@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('James', 'Wilson', 'jameswilson', 'password123', 'james.wilson@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Linda', 'Moore', 'lindamoore', 'password123', 'linda.moore@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Robert', 'Taylor', 'roberttaylor', 'password123', 'robert.taylor@example.com');
INSERT INTO public.users (first_name, last_name, username, password, email) VALUES ('Patricia', 'Anderson', 'patriciaanderson', 'password123', 'patricia.anderson@example.com');

-- Insert queries for 'blogs' table with reference to 'users' table
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('First Blog Post', 'Content of the first blog post', 1, 'first-blog-post', 'This is the excerpt of the first blog post', 'published');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Second Blog Post', 'Content of the second blog post', 2, 'second-blog-post', 'This is the excerpt of the second blog post', 'published');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Third Blog Post', 'Content of the third blog post', 3, 'third-blog-post', 'This is the excerpt of the third blog post', 'draft');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Fourth Blog Post', 'Content of the fourth blog post', 4, 'fourth-blog-post', 'This is the excerpt of the fourth blog post', 'draft');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Fifth Blog Post', 'Content of the fifth blog post', 5, 'fifth-blog-post', 'This is the excerpt of the fifth blog post', 'review');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Sixth Blog Post', 'Content of the sixth blog post', 6, 'sixth-blog-post', 'This is the excerpt of the sixth blog post', 'review');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Seventh Blog Post', 'Content of the seventh blog post', 7, 'seventh-blog-post', 'This is the excerpt of the seventh blog post', 'published');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Eighth Blog Post', 'Content of the eighth blog post', 8, 'eighth-blog-post', 'This is the excerpt of the eighth blog post', 'draft');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Ninth Blog Post', 'Content of the ninth blog post', 9, 'ninth-blog-post', 'This is the excerpt of the ninth blog post', 'review');
INSERT INTO public.blogs (title, content, author_id, slug, excerpt, status) VALUES ('Tenth Blog Post', 'Content of the tenth blog post', 10, 'tenth-blog-post', 'This is the excerpt of the tenth blog post', 'published');




CREATE TABLE public.categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE public.blog_categories (
    blog_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (blog_id, category_id),
    FOREIGN KEY (blog_id) REFERENCES public.blogs(id),
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
    blog_id INT NOT NULL,
    image_url TEXT NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (blog_id) REFERENCES public.blogs(id)
);

-- Likes table
CREATE TABLE public.likes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    blog_id INT NOT NULL REFERENCES blogs(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Comments table (simplified for hierarchical modeling)
CREATE TABLE public.comments (
    id SERIAL PRIMARY KEY,
    blog_id INT NOT NULL REFERENCES blogs(id),
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
    blog_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (blog_id, tag_id),
    FOREIGN KEY (blog_id) REFERENCES public.blogs(id),
    FOREIGN KEY (tag_id) REFERENCES public.tags(id)
);


-- Indexes for optimization
CREATE INDEX idx_users_username ON public.users(username);
CREATE INDEX idx_blogs_author_id ON public.blogs(author_id);
CREATE INDEX idx_images_blog_id ON public.images(blog_id);
CREATE INDEX idx_blogs_slug ON public.blogs(slug);
CREATE INDEX idx_blogs_published_at ON public.blogs(published_at);

ALTER TABLE public.blogs
DROP CONSTRAINT IF EXISTS blogs_author_id_fkey,
ADD CONSTRAINT blogs_author_id_fkey FOREIGN KEY (author_id)
REFERENCES public.users (id) ON DELETE CASCADE;


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

-- Trigger for blogs table
CREATE TRIGGER update_blog_modtime
BEFORE UPDATE ON public.blogs
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