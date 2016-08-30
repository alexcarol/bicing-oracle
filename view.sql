DROP VIEW fit_precalculation;

CREATE VIEW fit_precalculation
AS SELECT s.id, s.updatetime, s.bikes >= 1 as pbikes, s.slots >= 1 as pslots, w.type as weather_type, w.temperature
    FROM station_state s, weather w
    WHERE w.time = (SELECT w2.time
    FROM weather w2
    WHERE w2.time >= s.updatetime
    ORDER BY w2.time ASC
    LIMIT 1
);
DROP VIEW fit_precalculation_2;

CREATE VIEW fit_precalculation_2
AS SELECT s.id, s.bikes, s.slots, s.updatetime, s.bikes >= 1 as pbikes, s.slots >= 1 as pslots, w.type as weather_type, w.temperature
    FROM station_state s, weather w
    WHERE w.time = (SELECT w2.time
    FROM weather w2
    WHERE w2.time >= s.updatetime
    ORDER BY w2.time ASC
    LIMIT 1
);
