BEGIN;

CREATE TABLE IF NOT EXISTS entity_definition (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  constraint_custom TEXT[]
);

CREATE TYPE attribute_type AS ENUM ('string', 'number', 'boolean', 'date', 'time', 'datetime', 'month', 'co-ordinate', 'file');

CREATE TABLE IF NOT EXISTS entity_attribute (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  list BOOLEAN,
  type ATTRIBUTE_TYPE NOT NULL,
  entity_id UUID NOT NULL,
  constraint_required BOOLEAN NOT NULL,
  constraint_unique BOOLEAN NOT NULL,
  constraint_min NUMERIC,
  constraint_max NUMERIC,
  constraint_pattern TEXT,
  constraint_enum TEXT[],
  constraint_custom TEXT[],
  FOREIGN KEY (entity_id) REFERENCES entity_definition(id) ON DELETE CASCADE
);

CREATE TYPE entity_relationship_cardinality AS ENUM ('one-to-one', 'one-to-many', 'many-to-one', 'many-to-many');

CREATE TABLE IF NOT EXISTS entity_relationship (
  id UUID PRIMARY KEY,
  cardinality ENTITY_RELATIONSHIP_CARDINALITY NOT NULL,
  source_entity_id UUID NOT NULL,
  source_attribute_id UUID NOT NULL,
  target_entity_id UUID NOT NULL,
  target_attribute_id UUID NOT NULL,
  FOREIGN KEY (source_entity_id) REFERENCES entity_definition(id),
  FOREIGN KEY (source_attribute_id) REFERENCES entity_attribute(id),
  FOREIGN KEY (target_entity_id) REFERENCES entity_definition(id),
  FOREIGN KEY (target_attribute_id) REFERENCES entity_attribute(id)
);

COMMIT;