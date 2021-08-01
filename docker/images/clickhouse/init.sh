#!/bin/bash
set -e

clickhouse client -n <<-EOSQL
  CREATE DATABASE IF NOT EXISTS eventhouse;

  CREATE TABLE IF NOT EXISTS eventhouse.events (
      EventId UInt64,
      PreviousEventId UInt64,
      EventType String,
      IdempotencyKey String,
      EntityType String,
      EntityId String,
      Recorded DateTime64(9),
      Payload String,
      INDEX IdempotencyKeyIdx IdempotencyKey TYPE minmax GRANULARITY 512,
      INDEX EventIdIdx EventId TYPE minmax GRANULARITY 512,
      INDEX EventTypeIdx EventType TYPE minmax GRANULARITY 512
  )
  ENGINE = MergeTree()
  ORDER BY (EntityType, EntityId, EventId)
  PARTITION BY EntityType;

  CREATE TABLE IF NOT EXISTS eventhouse.apiClients (
      ClientId String,
      SecretHash String
  )
  ENGINE = MergeTree()
  ORDER BY (ClientId);
EOSQL