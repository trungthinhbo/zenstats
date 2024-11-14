CREATE TABLE visit (
  ip inet not null,
  visited_at timestamptz not null,
  primary key (ip, visited_at)
);

CREATE INDEX ON visit(visited_at);

CREATE INDEX ON visit(ip);
