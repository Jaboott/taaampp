import os

from flask import Flask, jsonify, request
from flask.cli import load_dotenv

from database import PostgresHandler

load_dotenv()
app = Flask(__name__)

def create_db_connection():
    return PostgresHandler(
        host=os.getenv('PG_HOST'),
        user=os.getenv('PG_USER'),
        password=os.getenv('PG_PASSWORD'),
        database=os.getenv('PG_DATABASE'),
        port=5432,
    )

@app.route('/')
def hello_world():  # put application's code here
    return 'Hello World!'

@app.route('/api/ping')
def ping_db():
    db = create_db_connection()
    try:
        db.execute("SELECT 1")

        return jsonify({
            'status': 'success',
            'message': 'Database connection successful',
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'fail',
            'message': str(e),
        }), 500
    finally:
        db.close()

@app.route('/api/register', methods=['POST'])
def register_user():
    from bcrypt import hashpw, gensalt
    data = request.get_json()

    if not data or not all(key in data for key in ['username', 'password', 'email']):
        return jsonify({
            'status': 'fail',
            'message': "Missing data. 'username', 'email', and 'password' are required.",
        }), 400

    username = data['username']
    hashed_password = hashpw(data['password'].encode('utf-8'), gensalt()).decode('utf-8')
    email = data['email']

    db = create_db_connection()
    try:
        db.execute(
            "INSERT INTO users (username, email, password_hash) VALUES (%s, %s, %s)",
            (username, email, hashed_password))
        return jsonify({
            'status': 'success',
            'message': 'User registered successfully',
        }), 201

    except Exception as e:
        if "unique constraint" in str(e).lower():
            return jsonify({
                'status': 'fail',
                'message': "Email or Username already exists. Please try another one.",
            }), 400
        return jsonify({
            'status': 'fail',
            'message': str(e),
        }), 500
    finally:
        db.close()


if __name__ == '__main__':
    app.run()
