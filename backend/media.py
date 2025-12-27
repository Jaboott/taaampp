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

@app.route('/api/media/<int:id>')
def get_media(id):
    db = create_db_connection()
    try:
        media_data = db.fetchone('select * from media where id = %s', (id,))
        return jsonify(media_data)
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500
    finally:
        db.close()

@app.route('/api/medias/<int:page>')
def get_medias(page):
    db = create_db_connection()
    if page == 0: page = 1
    try:
        medias = db.fetchall('select * from media order by id limit 50 offset %s', ((page-1)*50,))
        return jsonify(medias)
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        })
    finally:
        db.close()

@app.route('/api/media_details/<int:id>')
def get_media_details(id):
    db = create_db_connection()
    try:
        media_data = db.fetchone('select * from media m left join media_details md on  m.id = md.id where m.id = %s', (id,))
        return jsonify(media_data)
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        })
    finally:
        db.close()


if __name__ == '__main__':
    app.run(debug=True)