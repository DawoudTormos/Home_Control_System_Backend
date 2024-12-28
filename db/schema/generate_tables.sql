BEGIN;


CREATE TABLE IF NOT EXISTS public.api_authentication
(
    auth_id serial NOT NULL,
    user_id integer NOT NULL,
    api_key character varying(256) COLLATE pg_catalog."default" NOT NULL,
    device_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp without time zone,
    CONSTRAINT api_authentication_pkey PRIMARY KEY (auth_id),
    CONSTRAINT api_authentication_api_key_key UNIQUE (api_key)
);

CREATE TABLE IF NOT EXISTS public.camera_command_logs
(
    log_id serial NOT NULL,
    camera_id integer NOT NULL,
    command_id integer NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT camera_command_logs_pkey PRIMARY KEY (log_id)
);

CREATE TABLE IF NOT EXISTS public.camera_commands
(
    command_id serial NOT NULL,
    camera_id integer NOT NULL,
    command_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT camera_commands_pkey PRIMARY KEY (command_id)
);

CREATE TABLE IF NOT EXISTS public.camera_snaps
(
    snap_id serial NOT NULL,
    camera_id integer NOT NULL,
    image_url text COLLATE pg_catalog."default" NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT camera_snaps_pkey PRIMARY KEY (snap_id)
);

CREATE TABLE IF NOT EXISTS public.cameras
(
    id integer NOT NULL DEFAULT nextval('cameras_camera_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    room_id integer NOT NULL,
    value character varying COLLATE pg_catalog."default" NOT NULL,
    color bigint NOT NULL,
    index integer NOT NULL,
    CONSTRAINT cameras_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.general_command_logs
(
    log_id serial NOT NULL,
    command_id integer NOT NULL,
    executed_by_user_id integer NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT general_command_logs_pkey PRIMARY KEY (log_id)
);

CREATE TABLE IF NOT EXISTS public.general_commands
(
    command_id serial NOT NULL,
    command_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT general_commands_pkey PRIMARY KEY (command_id)
);

CREATE TABLE IF NOT EXISTS public.rooms
(
    id integer NOT NULL DEFAULT nextval('rooms_room_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    user_id integer NOT NULL,
    index integer NOT NULL DEFAULT 0,
    CONSTRAINT rooms_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.sensor_command_logs
(
    log_id serial NOT NULL,
    sensor_id integer NOT NULL,
    command_id integer NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sensor_command_logs_pkey PRIMARY KEY (log_id)
);

CREATE TABLE IF NOT EXISTS public.sensor_commands
(
    command_id serial NOT NULL,
    sensor_id integer NOT NULL,
    command_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT sensor_commands_pkey PRIMARY KEY (command_id)
);

CREATE TABLE IF NOT EXISTS public.sensor_logs
(
    log_id serial NOT NULL,
    sensor_id integer NOT NULL,
    reading double precision NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sensor_logs_pkey PRIMARY KEY (log_id)
);

CREATE TABLE IF NOT EXISTS public.sensor_types
(
    type_id serial NOT NULL,
    type_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT sensor_types_pkey PRIMARY KEY (type_id)
);

CREATE TABLE IF NOT EXISTS public.sensors
(
    id integer NOT NULL DEFAULT nextval('sensors_sensor_id_seq'::regclass),
    type character varying(50) COLLATE pg_catalog."default" NOT NULL,
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    room_id integer NOT NULL,
    type_id integer NOT NULL,
    value integer NOT NULL,
    index integer NOT NULL,
    color bigint NOT NULL,
    CONSTRAINT sensors_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.switch_command_logs
(
    log_id serial NOT NULL,
    switch_id integer NOT NULL,
    command_id integer NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT switch_command_logs_pkey PRIMARY KEY (log_id)
);

CREATE TABLE IF NOT EXISTS public.switch_commands
(
    command_id serial NOT NULL,
    switch_id integer NOT NULL,
    command_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT switch_commands_pkey PRIMARY KEY (command_id)
);

CREATE TABLE IF NOT EXISTS public.switch_logs
(
    log_id serial NOT NULL,
    switch_id integer NOT NULL,
    state boolean NOT NULL,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT switch_logs_pkey PRIMARY KEY (log_id)
);

CREATE TABLE IF NOT EXISTS public.switches
(
    id integer NOT NULL DEFAULT nextval('switches_switch_id_seq'::regclass),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    room_id integer NOT NULL,
    color bigint NOT NULL,
    icon_code integer NOT NULL,
    icon_family character varying COLLATE pg_catalog."default" NOT NULL,
    index integer NOT NULL,
    type integer NOT NULL,
    value smallint NOT NULL,
    CONSTRAINT switches_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.switches_schedule
(
    schedule_id serial NOT NULL,
    switch_id integer NOT NULL,
    command_id integer NOT NULL,
    scheduled_time timestamp without time zone NOT NULL,
    is_recurring boolean DEFAULT false,
    recurrence_interval interval,
    CONSTRAINT switches_schedule_pkey PRIMARY KEY (schedule_id)
);

CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    email character varying(150) COLLATE pg_catalog."default" NOT NULL,
    hashed_password text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email)
);

ALTER TABLE IF EXISTS public.api_authentication
    ADD CONSTRAINT api_authentication_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.camera_command_logs
    ADD CONSTRAINT camera_command_logs_camera_id_fkey FOREIGN KEY (camera_id)
    REFERENCES public.cameras (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.camera_command_logs
    ADD CONSTRAINT camera_command_logs_command_id_fkey FOREIGN KEY (command_id)
    REFERENCES public.camera_commands (command_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.camera_snaps
    ADD CONSTRAINT camera_snaps_camera_id_fkey FOREIGN KEY (camera_id)
    REFERENCES public.cameras (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.cameras
    ADD CONSTRAINT cameras_room_id_fkey FOREIGN KEY (room_id)
    REFERENCES public.rooms (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.general_command_logs
    ADD CONSTRAINT general_command_logs_command_id_fkey FOREIGN KEY (command_id)
    REFERENCES public.general_commands (command_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.general_command_logs
    ADD CONSTRAINT general_command_logs_user_id_fkey FOREIGN KEY (executed_by_user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.rooms
    ADD CONSTRAINT rooms_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES public.users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.sensor_command_logs
    ADD CONSTRAINT sensor_command_logs_command_id_fkey FOREIGN KEY (command_id)
    REFERENCES public.sensor_commands (command_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.sensor_command_logs
    ADD CONSTRAINT sensor_command_logs_sensor_id_fkey FOREIGN KEY (sensor_id)
    REFERENCES public.sensors (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.sensor_logs
    ADD CONSTRAINT sensor_logs_sensor_id_fkey FOREIGN KEY (sensor_id)
    REFERENCES public.sensors (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.sensors
    ADD CONSTRAINT sensors_room_id_fkey FOREIGN KEY (room_id)
    REFERENCES public.rooms (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.sensors
    ADD CONSTRAINT sensors_type_id_fkey FOREIGN KEY (type_id)
    REFERENCES public.sensor_types (type_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE SET NULL;


ALTER TABLE IF EXISTS public.switch_command_logs
    ADD CONSTRAINT switch_command_logs_command_id_fkey FOREIGN KEY (command_id)
    REFERENCES public.switch_commands (command_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.switch_command_logs
    ADD CONSTRAINT switch_command_logs_switch_id_fkey FOREIGN KEY (switch_id)
    REFERENCES public.switches (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.switch_logs
    ADD CONSTRAINT switch_logs_switch_id_fkey FOREIGN KEY (switch_id)
    REFERENCES public.switches (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.switches
    ADD CONSTRAINT switches_room_id_fkey FOREIGN KEY (room_id)
    REFERENCES public.rooms (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.switches_schedule
    ADD CONSTRAINT switches_schedule_command_id_fkey FOREIGN KEY (command_id)
    REFERENCES public.switch_commands (command_id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;


ALTER TABLE IF EXISTS public.switches_schedule
    ADD CONSTRAINT switches_schedule_switch_id_fkey FOREIGN KEY (switch_id)
    REFERENCES public.switches (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE;

END;