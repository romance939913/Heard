from flask import Flask, request, jsonify, make_response
from flask_bcrypt import Bcrypt
from flask_jwt_extended import JWTManager, create_access_token, set_access_cookies, jwt_required, get_jwt_identity
from flask_sqlalchemy import SQLAlchemy
from flask_cors import CORS
from datetime import datetime, timedelta, timezone
from config import Config

# setup Flask app
app = Flask(__name__)
CORS(app)

# setup password encryption
bcrypt = Bcrypt(app)

# setup JWT and database
app.config['SQLALCHEMY_DATABASE_URI'] = f"postgresql://brennanromance:{Config.POSTGRESQL_PASSWORD}@localhost/heard"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.config['JWT_TOKEN_LOCATION'] = ['cookies']
app.config['JWT_ACCESS_COOKIE_PATH'] = '/api/'
app.config["JWT_SECRET_KEY"] = Config.SECRET_KEY_JWT  
app.config["JWT_ACCESS_TOKEN_EXPIRES"] = timedelta(hours=168)

db = SQLAlchemy(app)
jwt = JWTManager(app)

# Models
class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(150), unique=True, nullable=False)
    email = db.Column(db.String(150), unique=True, nullable=False)
    password = db.Column(db.String(500), nullable=False)

class Company(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(150), unique=True, nullable=False)
    headquarters = db.Column(db.String(150))
    industry = db.Column(db.String(150))

    def to_dict(self):
        return {
            'id': self.id,
            'name': self.name,
            'headquarters': self.headquarters,
            'industry': self.industry
        }

class Suggestion(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, nullable=False)
    company_id = db.Column(db.Integer, nullable=False)
    title = db.Column(db.String(150), nullable=False)
    description = db.Column(db.Text)
    product_id = db.Column(db.Integer)
    upvotes = db.Column(db.Integer, nullable=False)
    downvotes = db.Column(db.Integer, nullable=False)

    def to_dict(self):
        return {
            'id': self.id,
            'title': self.title,
            'description': self.description,
            'user_id': self.user_id,
            'company_id': self.company_id,
            'product_id': self.product_id,
            'upvotes': self.upvotes,
            'downvotes': self.downvotes
        }

class Product(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(150), nullable=False)
    company_id = db.Column(db.Integer, nullable=False)
    description = db.Column(db.Text)

    def to_dict(self):
        return {
            'id': self.id,
            'name': self.name,
            'description': self.description,
            'company_id': self.company_id,
        }
    

# Create the database tables (run this only once)
with app.app_context():
    db.create_all()


# company routes
@app.post("/api/company")
@jwt_required()
def create_company():
    data = request.get_json()
    name = data["name"]
    hq = data["headquarters"]
    industry = data["industry"]

    if not data or not data.get('name'):
        return jsonify({"msg": "Missing required fields"}), 400

    new_user = Company(name=name, headquarters=hq, industry=industry)
    db.session.add(new_user)
    db.session.commit()
    return jsonify(name=name, headquarters=hq, industry=industry), 201


@app.get("/api/company")
@jwt_required()
def get_companies():
    companies = Company.query.all()
    return jsonify([company.to_dict() for company in companies]), 201


# product routes
@app.post("/api/product")
@jwt_required()
def create_product():
    data = request.get_json()
    name = data["name"]
    company_id = data["company_id"]
    description = data["description"]

    if not data or not data.get('name') or not data.get('description') or not data.get('company_id'):
        return jsonify({"msg": "Missing required fields"}), 400

    new_product = Product(name=name, company_id=company_id, description=description)
    db.session.add(new_product)
    db.session.commit()
    return jsonify(name=name, company_id=company_id, description=description), 201


@app.get("/api/product")
@jwt_required()
def get_products():
    data = request.get_json()
    
    if not data or not data.get("company_id"):
        return jsonify({"msg": "Missing required fields: company_id"}), 400
    
    company_id = data["company_id"]
    
    products = Product.query.filter_by(company_id=company_id)
    return jsonify([product.to_dict() for product in products]), 201


# suggestion routes
@app.post("/api/suggestion")
@jwt_required()
def create_suggestion():
    data = request.get_json()
    
    if not data or not data.get("title") or not data.get("user_id") or not data.get("company_id") or not data.get("description") or not data.get("product_id"):
        return jsonify({"msg": "Missing required fields"}), 400
    
    title, description, user_id, company_id, product_id = data["title"], data["description"], data["user_id"], data["company_id"], data["product_id"]

    suggestion = Suggestion.query.filter_by(title=title).first()
    if suggestion:
        return jsonify({"msg": "another suggestion with that title already exists"}), 400

    new_suggestion = Suggestion(title=title, description=description, user_id=user_id, company_id=company_id, product_id=product_id, upvotes=0, downvotes=0)
    db.session.add(new_suggestion)
    db.session.commit()

    return jsonify(title=title, description=description, user_id=user_id, company_id=company_id, product_id=product_id, upvotes=0, downvotes=0), 201


@app.patch("/api/suggestion")
@jwt_required()
def update_suggestion():
    data = request.get_json()
    
    if not data or not data.get("id"):
        return jsonify({"msg": "Missing required fields"}), 400
    
    id = data["id"]
    suggestion = Suggestion.query.get(id)

    if 'upvotes' in data:
        suggestion.upvotes = data["upvotes"]
    if 'downvotes' in data:
        suggestion.downvotes = data["downvotes"]
    if 'title' in data:
        suggestion.title = data["title"]
    if 'description' in data:
        suggestion.description = data["description"]

    db.session.commit()
    return jsonify(suggestion.to_dict()), 201


@app.get("/api/suggestion")
@jwt_required()
def get_suggestions():
    data = request.get_json()

    if not data or not data.get("product_id_included") or not data.get("filter_id"):
        return jsonify({"msg": "Missing required fields"}), 400
    
    product_id_included = data["product_id_included"]
    filter_id = data['filter_id']

    suggestions = None
    if product_id_included:
        suggestions = Suggestion.query.filter_by(product_id=filter_id)
    else:
        suggestions = Suggestion.query.filter_by(company_id=filter_id)
    return jsonify([suggestion.to_dict() for suggestion in suggestions]), 201


@app.delete("/api/suggestion")
@jwt_required()
def delete_suggestion():
    data = request.get_json()

    if not data or not data.get("id") or not data.get("user_id"):
        return jsonify({"msg": "Missing required fields"}), 400
    
    user_id = data["user_id"]
    id = data["id"]

    suggestion = db.session.query(Suggestion).filter_by(id=id).first()
    if suggestion.user_id == user_id:
        db.session.delete(suggestion)
        db.session.commit()
        return jsonify(suggestion.to_dict()), 201
    else:
        return jsonify({"msg": "user_id doesn't match"}), 400


# user auth routes
@app.route("/api/register", methods=["POST"])
def register():
    data = request.get_json()

    if not data or not data.get('username') or not data.get('password'):
        return jsonify({"msg": "Missing required fields"}), 400

    username = data['username']
    email = data['email']
    password = data['password']

    user = User.query.filter_by(username=username).first()
    if user:
        return jsonify({"msg": "Username already exists"}), 400
    
    hashed_password = bcrypt.generate_password_hash(password).decode('utf-8')
    new_user = User(username=username, password=hashed_password, email=email)
    db.session.add(new_user)
    db.session.commit()

    access_token = create_access_token(identity=new_user.email)
    response = make_response("Cookie set")
    response.set_cookie('access_token', access_token, path='/api')
    return jsonify(access_token=access_token), 200


@app.route("/api/login", methods=["POST"])
def login():
    data = request.get_json()

    if not data or not data.get('email') or not data.get('password'):
        return jsonify({"msg": "Missing required fields"}), 400

    email = data['email']
    password = data['password']

    user = User.query.filter_by(email=email).first()

    if not user or not bcrypt.check_password_hash(user.password, password):
        return jsonify({"msg": "Invalid credentials"}), 401

    access_token = create_access_token(identity=user.email)
    response = make_response("Cookie set")
    response.set_cookie('access_token', access_token, path='/api')
    return jsonify(access_token=access_token), 200


# @app.after_request
# def refresh_expiring_jwts(response):
#     try:
#         exp_timestamp = get_jwt()["exp"]
#         now = datetime.now(timezone.utc)
#         target_timestamp = datetime.timestamp(now + timedelta(minutes=30))
#         if target_timestamp > exp_timestamp:
#             access_token = create_access_token(identity=get_jwt_identity())
#         return response
#     except (RuntimeError, KeyError):
#         # Case where there is not a valid JWT. Just return the original response
#         return response