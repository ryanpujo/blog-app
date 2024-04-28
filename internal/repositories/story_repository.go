package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
)

type StoryRepository interface {
	Create(blog models.StoryPayload) (*uint, error)
	FindById(id uint) (*models.Story, error)
	FindBlogs() ([]*models.Story, error)
	DeleteById(id uint) error
	Update(id uint, payload models.StoryPayload) error
}

type storyRepository struct {
	Db *sql.DB
}

func NewBlogRepository(db *sql.DB) *storyRepository {
	return &storyRepository{
		Db: db,
	}
}

// Create inserts a new blog entry into the blogs table.
// It returns the ID of the newly inserted blog post or an error if the operation fails.
func (repo *storyRepository) Create(blog models.StoryPayload) (*uint, error) {
	// Set a timeout for the database operation.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Prepare the SQL statement for inserting a new blog post.
	stmt := `
		INSERT INTO stories (title, content, author_id, slug, excerpt, type, word_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
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
		blog.Type.String(),
		blog.WordCount,
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
func (repo *storyRepository) FindById(id uint) (*models.Story, error) {
	// Create a context with a timeout to avoid long-running queries.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to select a blog and its author's details.
	stmt := `
	SELECT b.id, b.title, b.content, b.slug, b.excerpt, b.status, b.published_at, b.updated_at, b.type, b.word_count,
	       u.id AS author_id, u.first_name, u.last_name, u.username, u.email
	FROM public.stories AS b
	INNER JOIN public.users AS u ON b.author_id = u.id
	WHERE b.id = $1;
	`

	// Prepare a Blog model to hold the data.
	var blog models.Story

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
		&blog.Type,
		&blog.WordCount,
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
func (repo *storyRepository) FindBlogs() ([]*models.Story, error) {
	// Create a context with a timeout to ensure the query does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to select all blogs and their authors' details.
	stmt := `
	SELECT b.id, b.title, b.content, b.slug, b.excerpt, b.status, b.published_at, b.updated_at, b.type, b.word_count,
	       u.id AS author_id, u.first_name, u.last_name, u.username, u.email
	FROM public.stories AS b
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
	blogs := []*models.Story{}

	// Iterate over the rows in the result set.
	for rows.Next() {
		var blog models.Story
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
			&blog.Type,
			&blog.WordCount,
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
func (repo *storyRepository) DeleteById(id uint) error {
	// Create a context with a timeout to ensure the operation does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to delete a blog post by ID.
	stmt := `
		DELETE FROM public.stories WHERE id = $1;
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
func (repo *storyRepository) Update(id uint, payload models.StoryPayload) error {
	// Create a context with a timeout to ensure the operation does not run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SQL statement to update a blog post.
	stmt := `
	UPDATE public.stories
	SET
		title = $1,
		content = $2,
		slug = $3,
		excerpt = $4,
		type = $5,
		word_count = $6
	WHERE id = $7;
	`

	// Execute the update statement with the provided payload and ID.
	result, err := repo.Db.ExecContext(ctx, stmt,
		payload.Title,
		payload.Content,
		payload.Slug,
		payload.Excerpt,
		payload.Type,
		payload.WordCount,
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
