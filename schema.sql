CREATE TABLE IF NOT EXISTS scraped_site (
  post_url TEXT PRIMARY KEY,
  post_title TEXT,
  feed_url TEXT,
  site_title TEXT,
  site_description TEXT,
  created TEXT
);