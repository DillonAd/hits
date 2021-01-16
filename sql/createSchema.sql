CREATE TABLE IF NOT EXISTS tenant
(
    id      UUID,
    name    VARCHAR(100),
    created TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tenant_page
(
    tenant_id UUID REFERENCES tenant(id),
    name      VARCHAR(100),
    created   TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (tenant_id, name)
);

CREATE TABLE IF NOT EXISTS hit
(
    tenant_id  UUID,
    page_name  VARCHAR(100),
    event_time TIMESTAMP WITH TIME ZONE,
    footprint  VARCHAR(128)

    CONSTRAINT fk_tenant_page FOREIGN KEY (tenant_id, page_name) REFERENCES tenant_page(tenant_id, name)
);