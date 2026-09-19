ALTER TABLE assignment
  ADD COLUMN active_order_marker CHAR(36) NULL AFTER active_marker;

UPDATE assignment
SET active_order_marker = service_order_id
WHERE is_active = 1;

ALTER TABLE assignment
  ADD CONSTRAINT uq_assignment_active_order_marker UNIQUE (active_order_marker),
  ADD CONSTRAINT ck_assignment_active_order_marker CHECK (
    (is_active = 1 AND active_order_marker IS NOT NULL AND active_order_marker = service_order_id)
    OR (is_active = 0 AND active_order_marker IS NULL)
  );
