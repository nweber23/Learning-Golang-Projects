PROJECT: Personal Blog
===============================================================================

OBJECTIVE:
Build a personal blog web application where you can write, publish, and manage
articles. The blog will have a public guest section for readers and a private
admin section for article management.

===============================================================================

GUEST SECTION (Public):
- Home Page: Display list of all published articles with titles and dates
- Article Page: Display full article content with publication date

ADMIN SECTION (Protected):
- Dashboard: List all articles with options to add, edit, or delete
- Add Article Page: Form to create new articles (title, content, date)
- Edit Article Page: Form to update existing articles
- Authentication: Basic session-based login for admin access

===============================================================================

KEY REQUIREMENTS:

Storage:
  • Use filesystem to store articles (JSON or Markdown format)
  • Each article stored as separate file
  • Include: title, content, publication date

Backend:
  • Server-side HTML rendering (no API required)
  • Form submission handling
  • Session-based authentication for admin section
  • Basic auth implementation (hardcoded credentials acceptable)

Frontend:
  • HTML and CSS (no JavaScript required)
  • Template engine for rendering articles
  • Responsive design for public and admin pages

===============================================================================

LEARNING OUTCOMES:
  ✓ Server-side templating
  ✓ Filesystem operations and file I/O
  ✓ Form handling and validation
  ✓ Basic authentication and session management
  ✓ Server-side HTML rendering
  ✓ CRUD operations for articles

FUTURE ENHANCEMENTS:
  • Comments on articles
  • Article categories and tags
  • Search functionality
  • Database integration (replace filesystem)
  • Advanced security measures

===============================================================================
