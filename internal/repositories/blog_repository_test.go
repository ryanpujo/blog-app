package repositories_test

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var excerpt = "a shorter post"
var blogPayload = models.BlogPayload{
	Title:    "my blog post",
	Content:  "a verry looooooooooong post",
	Slug:     "my-blog-post",
	AuthorID: 1,
	Excerpt:  &excerpt,
}

func Test_blogRepo_Create(t *testing.T) {
	id, err := blogRepo.Create(blogPayload)
	require.NoError(t, err)
	require.Equal(t, uint(11), *id)

	id, err = blogRepo.Create(blogPayload)
	require.Nil(t, id)
	require.Error(t, err)
	var pgErr *pgconn.PgError
	if assert.ErrorAs(t, err, &pgErr) {
		require.Equal(t, utils.ErrCodeUniqueViolation, pgErr.Code)
	}
}

func Test_blogRepo_FindById(t *testing.T) {
	blog, err := blogRepo.FindById(5)
	require.NoError(t, err)
	require.NotNil(t, blog)
	require.Equal(t, "Fifth Blog Post", blog.Title)
	require.Equal(t, "davidjones", blog.Author.Username)

	blog, err = blogRepo.FindById(20)
	require.Error(t, err)
	require.Nil(t, blog)
	if assert.ErrorAs(t, err, &sql.ErrNoRows) {
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func Test_blogRepo_FindBlogs(t *testing.T) {
	blogs, err := blogRepo.FindBlogs()
	require.NoError(t, err)
	require.NotNil(t, blogs)
	require.Equal(t, 11, len(blogs))
}

func Test_blogRepo_DeleteById(t *testing.T) {
	err := blogRepo.DeleteById(4)
	require.NoError(t, err)

	blog, err := blogRepo.FindById(4)
	require.Nil(t, blog)
	require.Error(t, err)

	err = blogRepo.DeleteById(23)
	require.Error(t, err)
	if assert.ErrorAs(t, err, &sql.ErrNoRows) {
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func Test_blogRepo_Update(t *testing.T) {
	blog, err := blogRepo.FindById(7)
	require.NoError(t, err)
	require.Equal(t, "Seventh Blog Post", blog.Title)

	updatePayload := models.BlogPayload{
		Title:   "The Seventh",
		Content: blog.Content,
		Slug:    blog.Slug,
		Excerpt: blog.Excerpt,
	}

	err = blogRepo.Update(7, updatePayload)
	require.NoError(t, err)

	updated, err := blogRepo.FindById(7)
	require.NoError(t, err)
	require.Equal(t, updatePayload.Title, updated.Title)

	err = blogRepo.Update(30, updatePayload)
	require.Error(t, err)
	if assert.ErrorAs(t, err, &sql.ErrNoRows) {
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}
