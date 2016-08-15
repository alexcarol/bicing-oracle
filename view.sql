CREATE VIEW fit_precalculation
AS SELECT id, updatetime, bikes > 1 as pbikes, slots < 1 as pslots
   FROM station_state;
