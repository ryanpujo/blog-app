package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

type BlogRepository interface {
	Create(blog models.BlogPayload) (*uint, error)
	FindById(id uint) (*models.Blog, error)
	FindBlogs() ([]*models.Blog, error)
	DeleteById(id uint) error
	Update(id uint, payload models.BlogPayload) error
}

type blogRepository struct {
	Db *sql.DB
}

func NewBlogRepository(db *sql.DB) *blogRepository {
	return &blogRepository{
		Db: db,
	}
}

// Create inserts a new blog entry into the blogs table.
// It returns the ID of the newly inserted blog post or an error if the operation fails.
func (repo *blogRepository) Create(blog models.BlogPayload) (*uint, error) {
	// Set a timeout for the database operation.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Prepare the SQL statement for inserting a new blog post.
	stmt := `
		INSERT INTO blogs (title, content, author_id, slug, excerpt)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	// Initialize the variable to store the returned ID.
	var id uint
	// Execute the SQL statement and scan the returned ID into the id variable.
	err := repo.Db.QueryRowContext(ctx, stmt,
		blog.Title,
		blog.Content,
		blog.AuthorID,
		blog.Slug,
		blog.Excerpt,
	).Scan(&id)
	if err != nil {
		// Handle any errors that occurred during the query execution.
		return nil, utils.HandlePostgresError(err)
	}

	// Return the pointer to the ID of the newly created blog post.
	return &id, nil
}

// FindById retrieves a blog post by its ID, including the author's information.
// It returns a pointer to a Blog model and any error encountered.
func (repo *blogRepository) FindById(id uint) (*models.Blog, error) {
	// Create a context with a timeout to avoid long-running queries.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to select a blog and its author's details.
	stmt := `
	SELECT b.id, b.title, b.content, b.slug, b.excerpt, b.status, b.published_at, b.updated_at,
	       u.id AS author_id, u.first_name, u.last_name, u.username, u.email
	FROM public.blogs AS b
	INNER JOIN public.users AS u ON b.author_id = u.id
	WHERE b.id = $1;
	`

	// Prepare a Blog model to hold the data.
	var blog models.Blog

	// Execute the query with the provided ID.
	row := repo.Db.QueryRowContext(ctx, stmt, id)

	// Scan the result into the Blog model.
	if err := row.Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.Slug,
		&blog.Excerpt,
		&blog.Status,
		&blog.PublishedAt,
		&blog.UpdatedAt,
		&blog.Author.ID,
		&blog.Author.FirstName,
		&blog.Author.LastName,
		&blog.Author.Username,
		&blog.Author.Email,
	); err != nil {
		// Handle any errors during scanning.
		return nil, utils.HandlePostgresError(err)
	}

	// Return a pointer to the populated Blog model.
	return &blog, nil
}

// FindBlogs retrieves all blog posts along with their corresponding authors' information.
// It returns a slice of pointers to Blog models and any error encountered.
func (repo *blogRepository) FindBlogs() ([]*models.Blog, error) {
	// Create a context with a timeout to ensure the query does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to select all blogs and their authors' details.
	stmt := `
	SELECT b.id, b.title, b.content, b.slug, b.excerpt, b.status, b.published_at, b.updated_at,
	       u.id AS author_id, u.first_name, u.last_name, u.username, u.email
	FROM public.blogs AS b
	INNER JOIN public.users AS u ON b.author_id = u.id
	`

	// Execute the query.
	rows, err := repo.Db.QueryContext(ctx, stmt)
	if err != nil {
		// Handle any errors that occur during query execution.
		return nil, utils.HandlePostgresError(err)
	}
	defer rows.Close()

	// Initialize a slice to hold the blog posts.
	blogs := []*models.Blog{}

	// Iterate over the rows in the result set.
	for rows.Next() {
		var blog models.Blog
		// Scan the result into the Blog model.
		if err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Slug,
			&blog.Excerpt,
			&blog.Status,
			&blog.PublishedAt,
			&blog.UpdatedAt,
			&blog.Author.ID,
			&blog.Author.FirstName,
			&blog.Author.LastName,
			&blog.Author.Username,
			&blog.Author.Email,
		); err != nil {
			// Handle any errors that occur during row scanning.
			return nil, utils.HandlePostgresError(err)
		}
		// Append the blog post to the slice.
		blogs = append(blogs, &blog)
	}

	// Check for any errors that might have occurred during row iteration.
	if err := rows.Err(); err != nil {
		return nil, utils.HandlePostgresError(err)
	}

	// Return the slice of blog posts.
	return blogs, nil
}

// DeleteById removes a blog post from the database by its ID.
// It returns an error if the deletion fails or if no record is found.
func (repo *blogRepository) DeleteById(id uint) error {
	// Create a context with a timeout to ensure the operation does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to delete a blog post by ID.
	stmt := `
		DELETE FROM public.blogs WHERE id = $1;
	`

	// Execute the delete statement.
	result, err := repo.Db.ExecContext(ctx, stmt, id)
	if err != nil {
		// Handle any errors that occur during the execution.
		return utils.HandlePostgresError(err)
	}

	// Check how many rows were affected by the delete operation.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Handle any errors that occur while checking the affected rows.
		return utils.HandlePostgresError(err)
	}

	// If no rows were affected, return an error indicating that no record was found.
	if rowsAffected == 0 {
		return utils.ErrNoDataFound
	}

	// Return nil if the deletion was successful.
	return nil
}

// Update modifies a blog post in the database using the provided ID and payload.
// It returns an error if the update operation fails or if no record is found.
func (repo *blogRepository) Update(id uint, payload models.BlogPayload) error {
	// Create a context with a timeout to ensure the operation does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to update a blog post.
	stmt := `
	UPDATE public.blogs
	SET
		title = $1,
		content = $2,
		slug = $3,
		excerpt = $4
	WHERE id = $5;
	`

	// Execute the update statement with the provided payload and ID.
	result, err := repo.Db.ExecContext(ctx, stmt,
		payload.Title,
		payload.Content,
		payload.Slug,
		payload.Excerpt,
		id,
	)
	if err != nil {
		// Handle any errors that occur during the execution.
		return utils.HandlePostgresError(err)
	}

	// Check how many rows were affected by the update operation.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Handle any errors that occur while checking the affected rows.
		return utils.HandlePostgresError(err)
	}

	// If no rows were affected, return an error indicating that no record was found.
	if rowsAffected == 0 {
		return utils.ErrNoDataFound
	}

	// Return nil if the update was successful.
	return nil
}
