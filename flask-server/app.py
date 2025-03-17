from flask import Flask, request, render_template, jsonify
from flask_bcrypt import Bcrypt
from flask_jwt_extended import JWTManager, create_access_token, jwt_required, get_jwt_identity
from flask_sqlalchemy import SQLAlchemy
from config import Config
import os
import psycopg2

# setup Flask app
app = Flask(__name__)

# setup password encryption
bcrypt = Bcrypt(app)

# setup JWT user authorization
app.config['SQLALCHEMY_DATABASE_URI'] = f"postgresql://postgres:{Config.POSTGRESQL_PASSWORD}@localhost/heard"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.config["JWT_SECRET_KEY"] = "super-secret"  # Change this in production

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

class Product(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(150), nullable=False)
    company_id = db.Column(db.Integer, nullable=False)
    description = db.Column(db.Text)


# Create the database tables (run this only once)
with app.app_context():
    db.create_all()

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


@app.post("/api/suggestion")
@jwt_required()
def create_suggestion():
    data = request.get_json()
    
    if not data or not data.get("title") or not data.get("user_id") or not data.get("company_id"):
        return jsonify({"msg": "Missing required fields"}), 400
    
    title, description, user_id, company_id, product_id = data["title"], data["description"], data["user_id"], data["company_id"], data["product_id"]
    if not product_id:
        product_id = None

    suggestion = Suggestion.query.filter_by(title=title).first()
    if suggestion:
        return jsonify({"msg": "another suggestion with that title already exists"}), 400

    new_suggestion = Suggestion(title=title, description=description, user_id=user_id, company_id=company_id, product_id=product_id, upvotes=0, downvotes=0)
    db.session.add(new_suggestion)
    db.session.commit(title=title, description=description, user_id=user_id, company_id=company_id, product_id=product_id, upvotes=0, downvotes=0)

    return jsonify({ "msg": "suggess" }), 200


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
    return jsonify(access_token=access_token), 200


        