package fmt

import (
	"strings"
	"testing"
)

func TestFormatSQL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		opts    SQLOptions
		want    string
		wantErr bool
	}{
		// SELECT with WHERE, ORDER BY, LIMIT.
		{
			name:  "basic_select",
			input: "select id, name from users where id = 1",
			opts:  SQLOptions{},
			want:  "select id, name\nfrom users\nwhere id = 1",
		},
		{
			name:  "select_with_where_and_or_orderby_limit",
			input: "SELECT id, name FROM users WHERE active = 1 AND role = 'admin' ORDER BY name LIMIT 10",
			opts:  SQLOptions{},
			want:  "select id, name\nfrom users\nwhere active = 1\n  and role = 'admin'\norder by name\nlimit 10",
		},
		{
			name:  "select_with_group_by",
			input: "SELECT department, COUNT(*) FROM employees GROUP BY department HAVING COUNT(*) > 5",
			opts:  SQLOptions{},
			want:  "select department, COUNT(*)\nfrom employees\ngroup by department\nhaving COUNT(*) > 5",
		},

		// JOIN.
		{
			name:  "inner_join",
			input: "select u.id, o.total from users u inner join orders o on u.id = o.user_id where o.total > 100",
			opts:  SQLOptions{Uppercase: true},
			want:  "SELECT u.id, o.total\nFROM users u\nINNER JOIN orders o\n  ON u.id = o.user_id\nWHERE o.total > 100",
		},
		{
			name:  "left_join",
			input: "select * from users left join orders on users.id = orders.user_id",
			opts:  SQLOptions{},
			want:  "select *\nfrom users\nleft join orders\n  on users.id = orders.user_id",
		},

		// INSERT.
		{
			name:  "insert_with_values",
			input: "insert into users (name, email) values ('John', 'john@example.com')",
			opts:  SQLOptions{Uppercase: true},
			want:  "INSERT INTO users (name, email)\nVALUES ('John', 'john@example.com')",
		},

		// UPDATE.
		{
			name:  "update_with_set_where",
			input: "update users set name = 'Jane' where id = 1",
			opts:  SQLOptions{Uppercase: true},
			want:  "UPDATE users\nSET name = 'Jane'\nWHERE id = 1",
		},

		// DELETE.
		{
			name:  "delete_from",
			input: "delete from users where id = 1",
			opts:  SQLOptions{Uppercase: true},
			want:  "DELETE FROM users\nWHERE id = 1",
		},

		// Uppercase flag.
		{
			name:  "uppercase_converts_keywords",
			input: "select id from users where id = 1",
			opts:  SQLOptions{Uppercase: true},
			want:  "SELECT id\nFROM users\nWHERE id = 1",
		},
		{
			name:  "lowercase_keywords_when_not_uppercase",
			input: "SELECT id FROM users WHERE id = 1",
			opts:  SQLOptions{Uppercase: false},
			want:  "select id\nfrom users\nwhere id = 1",
		},

		// UNION.
		{
			name:  "union",
			input: "select id from users union select id from admins",
			opts:  SQLOptions{Uppercase: true},
			want:  "SELECT id\nFROM users\nUNION\nSELECT id\nFROM admins",
		},

		// Subquery (formatted as part of the text between keywords).
		{
			name:  "subquery_in_where",
			input: "select * from users where id in (select user_id from orders)",
			opts:  SQLOptions{},
		},

		// Extra whitespace normalization.
		{
			name:  "extra_whitespace",
			input: "  select   id  from   users   where  id = 1  ",
			opts:  SQLOptions{},
			want:  "select id\nfrom users\nwhere id = 1",
		},

		// AND/OR as sub-clauses.
		{
			name:  "multiple_and_or",
			input: "select * from t where a = 1 and b = 2 or c = 3",
			opts:  SQLOptions{},
			want:  "select *\nfrom t\nwhere a = 1\n  and b = 2\n  or c = 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatSQL([]byte(tt.input), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			gotStr := strings.TrimSpace(string(got))
			if tt.want != "" {
				wantStr := strings.TrimSpace(tt.want)
				if gotStr != wantStr {
					t.Errorf("FormatSQL() =\n%s\nwant:\n%s", gotStr, wantStr)
				}
			}
			// For cases without specific want, just verify no error.
		})
	}
}

func TestFormatSQL_PreservesIdentifiers(t *testing.T) {
	// Uppercase option should convert keywords but preserve table/column names.
	input := "select username, email_address from user_accounts where is_active = true"
	got, err := FormatSQL([]byte(input), SQLOptions{Uppercase: true})
	if err != nil {
		t.Fatalf("FormatSQL: %v", err)
	}
	result := string(got)
	if !strings.Contains(result, "username") {
		t.Error("identifier 'username' was modified")
	}
	if !strings.Contains(result, "email_address") {
		t.Error("identifier 'email_address' was modified")
	}
	if !strings.Contains(result, "user_accounts") {
		t.Error("identifier 'user_accounts' was modified")
	}
	if !strings.Contains(result, "SELECT") {
		t.Error("keyword SELECT not uppercased")
	}
}
