-- Drop constraints first
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

-- Drop indexes
DROP INDEX IF EXISTS "accounts_owner_currency_idx";

-- Drop table
DROP TABLE IF EXISTS "users";
