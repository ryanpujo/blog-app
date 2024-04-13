package repositories_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanpujo/blog-app/models"
	"github.com/ryanpujo/blog-app/utils"
	"github.com/stretchr/testify/require"
)

var excerpt = "a shorter post"
var blogPayload = models.BlogPayload{
	Title:    "my blog post",
	Content:  "a very long post",
	Slug:     "my-blog-post",
	AuthorID: 1,
	Excerpt:  &excerpt,
}
var id = uint(1)
var expectExcerpt = "Test excerpt"
var updatedAt = time.Now()
var expectedBlog = &models.Blog{
	ID:          id,
	Title:       "Test Blog",
	Content:     "This is a test blog content.",
	Slug:        "test-blog",
	Excerpt:     &expectExcerpt,
	Status:      "published",
	PublishedAt: &updatedAt,
	UpdatedAt:   &updatedAt,
	Author: models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john.doe@example.com",
	},
}

// Test_blogRepo_Create tests the Create method of the blog repository.
func Test_blogRepo_Create(t *testing.T) {
	// Define a table-driven test with different scenarios.
	testTable := map[string]struct {
		blog    models.BlogPayload
		arrange func(mock sqlmock.Sqlmock)
		assert  func(t *testing.T, actualID *uint, err error)
	}{
		// Test case for successful blog creation.
		"success": {
			blog: blogPayload,
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO blogs").
					WithArgs("my blog post", "a very long post", 1, "my-blog-post", "a shorter post").
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.NoError(t, err)
				require.NotNil(t, actualID)
				require.Equal(t, uint(1), *actualID)
			},
		},
		// Test case for failure due to scanning error.
		"failed to scan": {
			blog: blogPayload,
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO blogs").
					WithArgs("my blog post", "a very long post", 1, "my-blog-post", "a shorter post").
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualID *uint, err error) {
				require.Error(t, err)
				require.Nil(t, actualID)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	// Iterate over each test case.
	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange(mock)

			id, err := blogRepo.Create(tc.blog)
			tc.assert(t, id, err)
		})
	}
}

// Test_blogRepo_FindById tests the FindById method of the blog repository.
func Test_blogRepo_FindById(t *testing.T) {

	// Define a table-driven test with different scenarios.
	testTable := map[string]struct {
		arrange func(mock sqlmock.Sqlmock)
		assert  func(t *testing.T, actualBlog *models.Blog, err error)
	}{
		// Test case for successful blog retrieval.
		"success": {
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"}).
					AddRow(expectedBlog.ID, expectedBlog.Title, expectedBlog.Content, expectedBlog.Slug,
						expectedBlog.Excerpt, expectedBlog.Status, expectedBlog.PublishedAt,
						expectedBlog.UpdatedAt, expectedBlog.Author.ID, expectedBlog.Author.FirstName,
						expectedBlog.Author.LastName, expectedBlog.Author.Username, expectedBlog.Author.Email)
				mock.ExpectQuery(`SELECT (.+) FROM public.blogs AS b INNER JOIN public.users AS u ON b.author_id = u.id`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualBlog *models.Blog, err error) {
				require.NoError(t, err)
				require.NotNil(t, actualBlog)
				require.Equal(t, expectedBlog, actualBlog)
			},
		},
		"failed": {
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})
				mock.ExpectQuery(`SELECT (.+) FROM public.blogs AS b INNER JOIN public.users AS u ON b.author_id = u.id`).
					WithArgs(id).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualBlog *models.Blog, err error) {
				require.Error(t, err)
				require.Nil(t, actualBlog)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	// Iterate over each test case.
	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange(mock)
			blog, err := blogRepo.FindById(id)
			tc.assert(t, blog, err)
		})
	}
}

func Test_blogRepo_FindBlogs(t *testing.T) {
	expectedBlogs := []*models.Blog{
		expectedBlog,
		expectedBlog,
	}
	testTable := map[string]struct {
		arrange func(mock sqlmock.Sqlmock)
		assert  func(t *testing.T, actualBlogs []*models.Blog, err error)
	}{
		// Test case for successful blog retrieval.
		"success": {
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})

				for _, expectedBlog := range expectedBlogs {
					rows.AddRow(expectedBlog.ID, expectedBlog.Title, expectedBlog.Content, expectedBlog.Slug,
						expectedBlog.Excerpt, expectedBlog.Status, expectedBlog.PublishedAt,
						expectedBlog.UpdatedAt, expectedBlog.Author.ID, expectedBlog.Author.FirstName,
						expectedBlog.Author.LastName, expectedBlog.Author.Username, expectedBlog.Author.Email)
				}

				mock.ExpectQuery(`SELECT (.+) FROM public.blogs AS b INNER JOIN public.users AS u ON b.author_id = u.id`).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualBlogs []*models.Blog, err error) {
				require.NoError(t, err)
				require.NotNil(t, actualBlogs)
				require.Equal(t, 2, len(actualBlogs))
			},
		},
		"failed": {
			arrange: func(mock sqlmock.Sqlmock) {
				// rows := sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})
				mock.ExpectQuery(`SELECT (.+) FROM public.blogs AS b INNER JOIN public.users AS u ON b.author_id = u.id`).
					WillReturnError(utils.ErrNoDataFound)
			},
			assert: func(t *testing.T, actualBlogs []*models.Blog, err error) {
				require.Error(t, err)
				require.Nil(t, actualBlogs)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"scan error": {
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"}).
					AddRow(expectedBlog.ID, expectedBlog.Title, expectedBlog.Content, expectedBlog.Slug,
						expectedBlog.Excerpt, expectedBlog.Status, expectedBlog.PublishedAt,
						expectedBlog.UpdatedAt, "expectedBlog.Author.ID", expectedBlog.Author.FirstName,
						expectedBlog.Author.LastName, expectedBlog.Author.Username, expectedBlog.Author.Email)
				mock.ExpectQuery(`SELECT (.+) FROM public.blogs AS b INNER JOIN public.users AS u ON b.author_id = u.id`).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualBlogs []*models.Blog, err error) {
				require.Error(t, err)
				require.Nil(t, actualBlogs)
			},
		},
		"row error": {
			arrange: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"}).
					AddRow(expectedBlog.ID, expectedBlog.Title, expectedBlog.Content, expectedBlog.Slug,
						expectedBlog.Excerpt, expectedBlog.Status, expectedBlog.PublishedAt,
						expectedBlog.UpdatedAt, expectedBlog.Author.ID, expectedBlog.Author.FirstName,
						expectedBlog.Author.LastName, expectedBlog.Author.Username, expectedBlog.Author.Email).RowError(0, utils.ErrNoDataFound)
				mock.ExpectQuery(`SELECT (.+) FROM public.blogs AS b INNER JOIN public.users AS u ON b.author_id = u.id`).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, actualBlogs []*models.Blog, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	// Iterate over each test case.
	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange(mock)
			blogs, err := blogRepo.FindBlogs()
			tc.assert(t, blogs, err)
		})
	}
}

func Test_blogRepo_DeleteById(t *testing.T) {
	testTable := map[string]struct {
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM public.blogs").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM public.blogs").WithArgs(1).WillReturnError(utils.ErrNoDataFound)
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"no record found": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM public.blogs").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"result error": {
			arrange: func() {
				sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectExec("DELETE FROM public.blogs").WithArgs(1).WillReturnResult(sqlmock.NewErrorResult(utils.ErrNoDataFound))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := blogRepo.DeleteById(1)

			tc.assert(t, err)
		})
	}
}

func Test_blogRepo_Update(t *testing.T) {
	testTable := map[string]struct {
		payload models.BlogPayload
		arrange func()
		assert  func(t *testing.T, err error)
	}{
		"success": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})

				mock.ExpectExec("UPDATE public.blogs SET").WithArgs(
					blogPayload.Title,
					blogPayload.Content,
					blogPayload.Slug,
					blogPayload.Excerpt,
					id,
				).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			assert: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		"failed": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})

				mock.ExpectExec("UPDATE public.blogs SET").WithArgs(
					blogPayload.Title,
					blogPayload.Content,
					blogPayload.Slug,
					blogPayload.Excerpt,
					id,
				).WillReturnError(utils.ErrNoDataFound)
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"no record Found": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})

				mock.ExpectExec("UPDATE public.blogs SET").WithArgs(
					blogPayload.Title,
					blogPayload.Content,
					blogPayload.Slug,
					blogPayload.Excerpt,
					id,
				).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
		"result error": {
			arrange: func() {
				sqlmock.NewRows([]string{"id", "title", "content", "slug", "excerpt", "status", "published_at", "updated_at", "author_id", "first_name", "last_name", "username", "email"})

				mock.ExpectExec("UPDATE public.blogs SET").WithArgs(
					blogPayload.Title,
					blogPayload.Content,
					blogPayload.Slug,
					blogPayload.Excerpt,
					id,
				).WillReturnResult(sqlmock.NewErrorResult(utils.ErrNoDataFound))
			},
			assert: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, utils.ErrNoDataFound, err)
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			tc.arrange()

			err := blogRepo.Update(id, blogPayload)

			tc.assert(t, err)
		})
	}
}
