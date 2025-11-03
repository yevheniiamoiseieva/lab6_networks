CREATE TABLE IF NOT EXISTS rsa_keys (
                                        id SERIAL PRIMARY KEY,
                                        public_key TEXT NOT NULL,
                                        private_key TEXT NOT NULL,
                                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);