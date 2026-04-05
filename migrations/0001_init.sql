CREATE TABLE IF NOT EXISTS manifest_verifications (
  id UUID PRIMARY KEY,
  manifest_id UUID NOT NULL,
  intent_id UUID NULL,
  pipeline_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  external_ref TEXT NOT NULL DEFAULT '',
  summary TEXT NOT NULL DEFAULT '',
  last_message TEXT NOT NULL DEFAULT '',
  steps JSONB NOT NULL DEFAULT '[]'::jsonb,
  details JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_observed_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_manifest_verifications_manifest_id
  ON manifest_verifications (manifest_id);

CREATE INDEX IF NOT EXISTS idx_manifest_verifications_intent_id
  ON manifest_verifications (intent_id);

CREATE INDEX IF NOT EXISTS idx_manifest_verifications_pipeline_id
  ON manifest_verifications (pipeline_id);

CREATE TABLE IF NOT EXISTS release_verifications (
  id UUID PRIMARY KEY,
  release_id UUID NOT NULL,
  intent_id UUID NULL,
  env TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  external_ref TEXT NOT NULL DEFAULT '',
  summary TEXT NOT NULL DEFAULT '',
  last_message TEXT NOT NULL DEFAULT '',
  steps JSONB NOT NULL DEFAULT '[]'::jsonb,
  details JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_observed_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_release_verifications_release_id
  ON release_verifications (release_id);

CREATE INDEX IF NOT EXISTS idx_release_verifications_intent_id
  ON release_verifications (intent_id);
