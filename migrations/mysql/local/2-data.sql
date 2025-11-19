USE db_bemax_api;

-- add roles default
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

-- add default categories to system
INSERT INTO REMINDER_CATEGORIES (id, user_id, name, name_key, description, icon, color, scope, display_order) VALUES
    (UUID(), NULL, 'Medication', 'category.medication', 'Medicine and prescriptions', '💊', '#FF5733', 'system', 1),
    (UUID(), NULL, 'Appointment', 'category.appointment', 'Medical appointments', '🩺', '#3498DB', 'system', 2),
    (UUID(), NULL, 'Exam', 'category.exam', 'Medical exams and tests', '🔬', '#9B59B6', 'system', 3),
    (UUID(), NULL, 'Exercise', 'category.exercise', 'Physical activities', '🏃', '#2ECC71', 'system', 4),
    (UUID(), NULL, 'Therapy', 'category.therapy', 'Therapy sessions', '🧠', '#E74C3C', 'system', 5),
    (UUID(), NULL, 'Diet', 'category.diet', 'Dietary reminders', '🥗', '#F39C12', 'system', 6),
    (UUID(), NULL, 'Other', 'category.other', 'Other reminders', '📌', '#95A5A6', 'system', 999);
