USE db_bemax_api;

-- Add roles default
INSERT INTO ROLES (id, name, description) VALUES
    (UUID(), 'admin', 'Administrador do sistema com acesso total'),
    (UUID(), 'manager', 'Gerente com acesso ao painel administrativo'),
    (UUID(), 'employee', 'Funcionário com acesso limitado'),
    (UUID(), 'customer', 'Cliente do sistema');

-- add all states Brazilian
INSERT INTO STATES (id, name, region) VALUES
    ('AC', 'Acre', 'Norte'),
    ('AL', 'Alagoas', 'Nordeste'),
    ('AM', 'Amazonas', 'Norte'),
    ('AP', 'Amapá', 'Norte'),
    ('BA', 'Bahia', 'Nordeste'),
    ('CE', 'Ceará', 'Nordeste'),
    ('DF', 'Distrito Federal', 'Centro-Oeste'),
    ('ES', 'Espírito Santo', 'Sudeste'),
    ('GO', 'Goiás', 'Centro-Oeste'),
    ('MA', 'Maranhão', 'Nordeste'),
    ('MG', 'Minas Gerais', 'Sudeste'),
    ('MS', 'Mato Grosso do Sul', 'Centro-Oeste'),
    ('MT', 'Mato Grosso', 'Centro-Oeste'),
    ('PA', 'Pará', 'Norte'),
    ('PB', 'Paraíba', 'Nordeste'),
    ('PE', 'Pernambuco', 'Nordeste'),
    ('PI', 'Piauí', 'Nordeste'),
    ('PR', 'Paraná', 'Sul'),
    ('RJ', 'Rio de Janeiro', 'Sudeste'),
    ('RN', 'Rio Grande do Norte', 'Nordeste'),
    ('RO', 'Rondônia', 'Norte'),
    ('RR', 'Roraima', 'Norte'),
    ('RS', 'Rio Grande do Sul', 'Sul'),
    ('SC', 'Santa Catarina', 'Sul'),
    ('SE', 'Sergipe', 'Nordeste'),
    ('SP', 'São Paulo', 'Sudeste'),
    ('TO', 'Tocantins', 'Norte');
