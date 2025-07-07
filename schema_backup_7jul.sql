--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13 (Debian 15.13-1.pgdg120+1)
-- Dumped by pg_dump version 15.13 (Debian 15.13-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: attendance; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.attendance (
    id integer NOT NULL,
    event_id integer,
    user_id integer,
    attendance_status character varying(10),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    token_input character varying(5),
    CONSTRAINT attendance_attendance_status_check CHECK (((attendance_status)::text = ANY ((ARRAY['present'::character varying, 'absent'::character varying])::text[])))
);


ALTER TABLE public.attendance OWNER TO admin;

--
-- Name: attendance_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.attendance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.attendance_id_seq OWNER TO admin;

--
-- Name: attendance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.attendance_id_seq OWNED BY public.attendance.id;


--
-- Name: certificates; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.certificates (
    id integer NOT NULL,
    event_id integer,
    user_id integer,
    certificate_data text NOT NULL,
    mint_status character varying(20),
    mint_transaction_hash character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    url_metadata character varying(512),
    url_certificate character varying(512),
    certificate_type character varying(100),
    CONSTRAINT certificates_mint_status_check CHECK (((mint_status)::text = ANY ((ARRAY['pending'::character varying, 'minted'::character varying, 'failed'::character varying])::text[])))
);


ALTER TABLE public.certificates OWNER TO admin;

--
-- Name: certificates_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.certificates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.certificates_id_seq OWNER TO admin;

--
-- Name: certificates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.certificates_id_seq OWNED BY public.certificates.id;


--
-- Name: event_certificates; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.event_certificates (
    event_id integer NOT NULL,
    url_certificate text NOT NULL,
    uploaded_by integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    description text
);


ALTER TABLE public.event_certificates OWNER TO admin;

--
-- Name: events; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.events (
    id integer NOT NULL,
    title character varying(255) NOT NULL,
    description text,
    vendor_id integer,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    status character varying(20) DEFAULT 'upcoming'::character varying,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    picture bytea NOT NULL,
    maxattendees integer DEFAULT 0 NOT NULL,
    location character varying(255) DEFAULT ''::character varying,
    requirements jsonb DEFAULT '[]'::jsonb,
    agenda jsonb DEFAULT '[]'::jsonb,
    token character varying(5) DEFAULT SUBSTRING(md5((random())::text) FROM 1 FOR 5) NOT NULL,
    CONSTRAINT events_status_check CHECK (((status)::text = ANY ((ARRAY['upcoming'::character varying, 'ongoing'::character varying, 'ended'::character varying, 'canceled'::character varying, 'minting'::character varying])::text[])))
);


ALTER TABLE public.events OWNER TO admin;

--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.events_id_seq OWNER TO admin;

--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    wallet_address character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO admin;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: vendors; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.vendors (
    id integer NOT NULL,
    vendor_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    contact_info text,
    wallet_address character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.vendors OWNER TO admin;

--
-- Name: vendors_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.vendors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.vendors_id_seq OWNER TO admin;

--
-- Name: vendors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.vendors_id_seq OWNED BY public.vendors.id;


--
-- Name: whitelist; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.whitelist (
    id integer NOT NULL,
    event_id integer,
    user_id integer,
    wallet_address character varying(255),
    status character varying(20),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT whitelist_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'approved'::character varying, 'rejected'::character varying])::text[])))
);


ALTER TABLE public.whitelist OWNER TO admin;

--
-- Name: whitelist_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.whitelist_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.whitelist_id_seq OWNER TO admin;

--
-- Name: whitelist_id_seq1; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.whitelist_id_seq1
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.whitelist_id_seq1 OWNER TO admin;

--
-- Name: whitelist_id_seq1; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.whitelist_id_seq1 OWNED BY public.whitelist.id;


--
-- Name: attendance id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attendance ALTER COLUMN id SET DEFAULT nextval('public.attendance_id_seq'::regclass);


--
-- Name: certificates id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.certificates ALTER COLUMN id SET DEFAULT nextval('public.certificates_id_seq'::regclass);


--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: vendors id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.vendors ALTER COLUMN id SET DEFAULT nextval('public.vendors_id_seq'::regclass);


--
-- Name: whitelist id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.whitelist ALTER COLUMN id SET DEFAULT nextval('public.whitelist_id_seq1'::regclass);


--
-- Data for Name: attendance; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.attendance (id, event_id, user_id, attendance_status, created_at, updated_at, token_input) FROM stdin;
1	3	3	present	2025-06-30 11:28:38.163372+07	2025-06-30 11:28:38.163372+07	G7N9L
2	4	3	present	2025-06-30 11:40:27.957797+07	2025-06-30 11:40:27.957797+07	7H83Z
\.


--
-- Data for Name: certificates; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.certificates (id, event_id, user_id, certificate_data, mint_status, mint_transaction_hash, created_at, updated_at, url_metadata, url_certificate, certificate_type) FROM stdin;
1	3	3	{"tokenURI":"ipfs://bafkreihjqmsyki64wh2oyr2vy2a247d5krqwxqb32xpamoymrh67l5zwde","urlMetadata":"https://bafkreihjqmsyki64wh2oyr2vy2a247d5krqwxqb32xpamoymrh67l5zwde.ipfs.w3s.link/","urlCertificate":"https://bafkreihatxutmfi2siq4okabgou7hm7ruu5nm2tks446ko45rwlghsuglq.ipfs.w3s.link/","certificateType":"Event Certificate Template for: Seminar Web 3","user_address":"0x467B8d4e4F3C6Ddf40F0C63092b500Ff4298cB6A"}	minted	0xe781297041aa51416ebc5a00a031c1bee509ab85d674a33d592c4174686ce497	2025-06-30 11:31:14.359868+07	2025-06-30 11:31:14.359868+07	https://bafkreihjqmsyki64wh2oyr2vy2a247d5krqwxqb32xpamoymrh67l5zwde.ipfs.w3s.link/	https://bafkreihatxutmfi2siq4okabgou7hm7ruu5nm2tks446ko45rwlghsuglq.ipfs.w3s.link/	Event Certificate Template for: Seminar Web 3
2	4	3	{"tokenURI":"ipfs://bafkreiat72ctw3o5mj432r6zcylhcwvptvmbyueotdd7cubraaayonm6pu","urlMetadata":"https://bafkreiat72ctw3o5mj432r6zcylhcwvptvmbyueotdd7cubraaayonm6pu.ipfs.w3s.link/","urlCertificate":"https://bafkreiafqovm3tqavvy4nbnsz7n7oerhkizmehqn23jgdkedsrwpdc62ei.ipfs.w3s.link/","certificateType":"Event Certificate Template for: Web 3 Seminar Development","user_address":"0x467B8d4e4F3C6Ddf40F0C63092b500Ff4298cB6A"}	minted	0x4a4460d848792e01f0cc471019ec3a3b4785d333776f42c86824a87d5eb16202	2025-06-30 11:41:41.031239+07	2025-06-30 11:41:41.031239+07	https://bafkreiat72ctw3o5mj432r6zcylhcwvptvmbyueotdd7cubraaayonm6pu.ipfs.w3s.link/	https://bafkreiafqovm3tqavvy4nbnsz7n7oerhkizmehqn23jgdkedsrwpdc62ei.ipfs.w3s.link/	Event Certificate Template for: Web 3 Seminar Development
\.


--
-- Data for Name: event_certificates; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.event_certificates (event_id, url_certificate, uploaded_by, created_at, description) FROM stdin;
3	https://bafkreihatxutmfi2siq4okabgou7hm7ruu5nm2tks446ko45rwlghsuglq.ipfs.w3s.link/	2	2025-06-30 11:30:49.983603	Event Certificate Template for: Seminar Web 3
4	https://bafkreiafqovm3tqavvy4nbnsz7n7oerhkizmehqn23jgdkedsrwpdc62ei.ipfs.w3s.link/	2	2025-06-30 11:40:54.137165	Event Certificate Template for: Web 3 Seminar Development
\.


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.events (id, title, description, vendor_id, start_date, end_date, status, created_at, updated_at, picture, maxattendees, location, requirements, agenda, token) FROM stdin;
1	Web3 Developer Workshop	Belajar Web3 dengan mentor dari Startup lokal Wong Liyo Ngerti Opo (WLNO)	1	2025-06-30 02:00:00+07	2025-06-30 03:00:00+07	upcoming	2025-06-30 03:08:57.911954+07	2025-06-30 03:08:57.911954+07	\\x75706c6f6164732f313735313232373733375f3539303739305f3635302e6a7067	30	Basement Gd. 5, Universitas Amikom Yogyakarta	["Minimal Mandi", "WIbu diperbolehkan", "Sekali lagi minimal mandi"]	[{"time": "08:30", "topic": "Open Gate"}, {"time": "09:10", "topic": "Opening"}, {"time": "09:30", "topic": "Workshop"}, {"time": "11:30", "topic": "Closing"}]	1A7XV
3	Seminar Web 3	Discuss abaout fundamental Web 3 	2	2025-06-30 03:00:00+07	2025-06-30 04:50:00+07	ongoing	2025-06-30 11:10:06.916977+07	2025-06-30 11:10:06.916977+07	\\x75706c6f6164732f313735313235363630365f5669727475616c204261636b67726f756e642042616e676b697420323032342e706e67	100	VR	["Nothing"]	[{"time": "13:10", "topic": "Start"}]	G7N9L
2	Seminar Nasional Informatika	Web 3 integration	2	2025-06-30 18:40:00+07	2025-06-30 19:00:00+07	upcoming	2025-06-30 10:33:02.68755+07	2025-06-30 10:33:02.68755+07	\\x75706c6f6164732f313735313235343338325f44656661756c745f44657369676e5f615f636c65616e5f616e645f6d6f6465726e5f6c6f676f5f666f725f615f7375737461696e61626c655f6f6e6c696e5f302e6a7067	10	Cinema Amikom	["Bawa alat tulis", "Laptop"]	[{"time": "11:30", "topic": "Opening"}, {"time": "12:00", "topic": "Closing"}]	F14PJ
4	Web 3 Seminar Development	Web 3 Seminar Development	2	2025-06-30 10:38:00+07	2025-06-30 11:38:00+07	upcoming	2025-06-30 11:39:14.139391+07	2025-06-30 11:39:14.139391+07	\\x75706c6f6164732f313735313235383335345f44656661756c745f44657369676e5f615f636c65616e5f616e645f6d6f6465726e5f6c6f676f5f666f725f615f7375737461696e61626c655f6f6e6c696e5f302e6a7067	100	VR	["Laptop"]	[{"time": "12:39", "topic": "bebas so fun"}]	7H83Z
5	Web 3 Sertifikasi	Web3	2	2025-07-04 02:51:00+07	2025-07-04 03:50:00+07	upcoming	2025-07-03 18:50:53.682017+07	2025-07-03 18:50:53.682017+07	\\x75706c6f6164732f313735313534333435335f646f776e6c6f6164202835292e6a706567	100	VR	["Laptop"]	[{"time": "20:00", "topic": "Opening"}]	J27SC
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (id, email, wallet_address, name, created_at, updated_at) FROM stdin;
1	gustipadaka19@gmail.com	0xaccB5E0993c482c95a0Cc4ed33958Fc689fa55D6	Gusti Padaka	2025-06-30 03:02:06.924534+07	2025-06-30 03:02:06.924534+07
3	ramasyailana3@gmail.com	0x467B8d4e4F3C6Ddf40F0C63092b500Ff4298cB6A	Rama Syailana Dewa	2025-06-30 11:06:51.920691+07	2025-06-30 11:06:51.920691+07
4	yantoalfurqan@gmail.com	0xc895Ac933134D018eaC9d53c2BD1eF2482e1515e	Antonio	2025-07-01 23:36:47.151962+07	2025-07-01 23:36:47.151962+07
\.


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.vendors (id, vendor_name, email, contact_info, wallet_address, created_at, updated_at) FROM stdin;
1	GPadaka Crop	vendor@gpadaka.com	08191919	0x6c0f62b1d1E20Af1f92088D4c0662a71bA5233a1	2025-06-30 03:03:53.58328+07	2025-06-30 03:03:53.58328+07
2	Amikom	amikom@amikom.ac.id		0x12d7A5E92D17dcb068e512660B24A9A3072a755e	2025-06-30 10:30:12.625795+07	2025-06-30 10:30:12.625795+07
\.


--
-- Data for Name: whitelist; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.whitelist (id, event_id, user_id, wallet_address, status, created_at, updated_at) FROM stdin;
1	2	1	0xaccB5E0993c482c95a0Cc4ed33958Fc689fa55D6	approved	2025-06-30 11:06:47.049733+07	2025-06-30 11:06:47.049733+07
2	3	3	0x467B8d4e4F3C6Ddf40F0C63092b500Ff4298cB6A	approved	2025-06-30 11:24:27.898507+07	2025-06-30 11:24:27.898507+07
3	4	3	0x467B8d4e4F3C6Ddf40F0C63092b500Ff4298cB6A	approved	2025-06-30 11:39:39.828935+07	2025-06-30 11:39:39.828935+07
\.


--
-- Name: attendance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.attendance_id_seq', 2, true);


--
-- Name: certificates_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.certificates_id_seq', 2, true);


--
-- Name: events_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.events_id_seq', 5, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_id_seq', 4, true);


--
-- Name: vendors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.vendors_id_seq', 2, true);


--
-- Name: whitelist_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.whitelist_id_seq', 1, false);


--
-- Name: whitelist_id_seq1; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.whitelist_id_seq1', 3, true);


--
-- Name: attendance attendance_event_user_unique; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attendance
    ADD CONSTRAINT attendance_event_user_unique UNIQUE (event_id, user_id);


--
-- Name: attendance attendance_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attendance
    ADD CONSTRAINT attendance_pkey PRIMARY KEY (id);


--
-- Name: certificates certificates_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.certificates
    ADD CONSTRAINT certificates_pkey PRIMARY KEY (id);


--
-- Name: event_certificates event_certificates_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.event_certificates
    ADD CONSTRAINT event_certificates_pkey PRIMARY KEY (event_id);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: events events_token_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_token_key UNIQUE (token);


--
-- Name: whitelist unique_event_user; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.whitelist
    ADD CONSTRAINT unique_event_user UNIQUE (event_id, user_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_wallet_address_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_wallet_address_key UNIQUE (wallet_address);


--
-- Name: vendors vendors_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_pkey PRIMARY KEY (id);


--
-- Name: vendors vendors_wallet_address_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_wallet_address_key UNIQUE (wallet_address);


--
-- Name: whitelist whitelist_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.whitelist
    ADD CONSTRAINT whitelist_pkey PRIMARY KEY (id);


--
-- Name: idx_whitelist_event_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_whitelist_event_id ON public.whitelist USING btree (event_id);


--
-- Name: idx_whitelist_event_status; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_whitelist_event_status ON public.whitelist USING btree (event_id, status);


--
-- Name: idx_whitelist_status; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_whitelist_status ON public.whitelist USING btree (status);


--
-- Name: idx_whitelist_user_id; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_whitelist_user_id ON public.whitelist USING btree (user_id);


--
-- Name: idx_whitelist_wallet_address; Type: INDEX; Schema: public; Owner: admin
--

CREATE INDEX idx_whitelist_wallet_address ON public.whitelist USING btree (wallet_address);


--
-- Name: attendance attendance_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attendance
    ADD CONSTRAINT attendance_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: attendance attendance_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attendance
    ADD CONSTRAINT attendance_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: certificates certificates_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.certificates
    ADD CONSTRAINT certificates_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: certificates certificates_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.certificates
    ADD CONSTRAINT certificates_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: event_certificates event_certificates_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.event_certificates
    ADD CONSTRAINT event_certificates_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: event_certificates event_certificates_uploaded_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.event_certificates
    ADD CONSTRAINT event_certificates_uploaded_by_fkey FOREIGN KEY (uploaded_by) REFERENCES public.vendors(id);


--
-- Name: events events_vendor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_vendor_id_fkey FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: whitelist whitelist_event_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.whitelist
    ADD CONSTRAINT whitelist_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: whitelist whitelist_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.whitelist
    ADD CONSTRAINT whitelist_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT USAGE ON SCHEMA public TO dewa;
GRANT USAGE ON SCHEMA public TO faruq;


--
-- Name: TABLE attendance; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.attendance TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.attendance TO faruq;


--
-- Name: TABLE certificates; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.certificates TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.certificates TO faruq;


--
-- Name: TABLE event_certificates; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.event_certificates TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.event_certificates TO faruq;


--
-- Name: TABLE events; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.events TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.events TO faruq;


--
-- Name: TABLE users; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.users TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.users TO faruq;


--
-- Name: TABLE vendors; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.vendors TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.vendors TO faruq;


--
-- Name: TABLE whitelist; Type: ACL; Schema: public; Owner: admin
--

GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.whitelist TO dewa;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.whitelist TO faruq;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: admin
--

ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public GRANT SELECT,INSERT,DELETE,UPDATE ON TABLES  TO dewa;
ALTER DEFAULT PRIVILEGES FOR ROLE admin IN SCHEMA public GRANT SELECT,INSERT,DELETE,UPDATE ON TABLES  TO faruq;


--
-- PostgreSQL database dump complete
--

ALTER TABLE public.vendors ADD COLUMN credit numeric(20,2) DEFAULT 0;
ALTER TABLE public.events ADD COLUMN event_fee numeric(20,2) DEFAULT 0;

   CREATE TABLE public.user_event_payments (
       id SERIAL PRIMARY KEY,
       user_id INTEGER NOT NULL REFERENCES public.users(id),
       event_id INTEGER NOT NULL REFERENCES public.events(id),
       amount_paid NUMERIC(20,2) NOT NULL,
       payment_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       status VARCHAR(20) DEFAULT 'pending', -- pending, paid, failed
       description TEXT
   );

   CREATE TABLE public.vendor_withdrawals (
       id SERIAL PRIMARY KEY,
       vendor_id INTEGER NOT NULL REFERENCES public.vendors(id),
       amount NUMERIC(20,2) NOT NULL,
       fee NUMERIC(20,2) NOT NULL,
       net_amount NUMERIC(20,2) NOT NULL,
       requested_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       status VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected, paid
       description TEXT
   );