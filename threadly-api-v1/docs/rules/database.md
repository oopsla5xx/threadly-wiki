# Database Rule: Naming (PostgreSQL)

Luôn đặt tên tường minh: `{table}_{columns}_{suffix}`. Không dùng tên tự sinh.


| Loại                          | Hậu tố   | Ví dụ                            |
| ----------------------------- | -------- | -------------------------------- |
| Primary key                   | `_pkey`  | `iam_users_pkey`                 |
| Unique constraint             | `_key`   | `iam_users_external_id_key`      |
| Unique index (gồm expression) | `_uidx`  | `iam_users_email_lower_uidx`     |
| Index thường                  | `_idx`   | `iam_users_created_at_idx`       |
| Foreign key                   | `_fkey`  | `iam_users_organization_id_fkey` |
| Check                         | `_check` | `iam_users_status_check`         |


- `UNIQUE (col)` → `_key`. `CREATE UNIQUE INDEX` (kể cả `LOWER(email)`) → `_uidx`, thêm tên biểu thức vào giữa.
- Tối đa 63 ký tự: rút gọn tên cột, giữ nguyên hậu tố.

```sql
CREATE UNIQUE INDEX iam_users_email_lower_uidx ON iam_users (LOWER(email));
CREATE INDEX iam_users_created_at_idx ON iam_users (created_at);

```