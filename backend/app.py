import os

from flask import Flask, jsonify
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
            'status': 'ok',
            'message': 'Database connection successful',
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'fail',
            'message': str(e),
        }), 500
    finally:
        db.close()


if __name__ == '__main__':
    app.run()
