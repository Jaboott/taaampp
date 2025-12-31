from datetime import datetime
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


def get_season():
    month = datetime.now().month

    if 1 <= month <= 3:
        return "WINTER"
    elif 4 <= month <= 6:
        return "SPRING"
    elif 7 <= month <= 9:
        return "SUMMER"
    else:
        return "FALL"


@app.route('/')
def hello_world():  # put application's code here
    return 'Hello World!'


@app.route('/api/media/<int:id>')
def get_media(id):
    db = create_db_connection()
    try:
        media_data = db.fetchone('select * from media where id = %s', (id,))
        return jsonify({
            'status': 'success',
            'data': media_data
        }), 200
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
        medias = db.fetchall('select * from media order by id limit 50 offset %s', ((page - 1) * 50,))
        return jsonify({
            'status': 'success',
            'data': medias
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500
    finally:
        db.close()


@app.route('/api/media_detail/<int:id>')
def get_media_details(id):
    db = create_db_connection()
    try:
        media_data = db.fetchone('select * from media m left join media_details md on  m.id = md.id where m.id = %s',
                                 (id,))
        return jsonify({
            'status': 'success',
            'data': media_data
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500
    finally:
        db.close()


@app.route('/api/popular_medias/<int:page>')
def get_popular(page):
    year = request.args.get('year')
    season = request.args.get('season')
    query_string = f'select m.* from media m left join media_details md on m.id = md.id where {f"m.season_year = {year}" if year else 'TRUE'} and {f"m.season = '{season}'" if season else 'TRUE'} order by md.popularity desc limit 50 offset %s'
    db = create_db_connection()
    try:
        popular_medias = db.fetchall(
            query_string,
            ((page - 1) * 50,))
        return jsonify({
            'status': 'success',
            'data': popular_medias
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500
    finally:
        db.close()


@app.route('/api/top_medias/<int:page>')
def get_top(page):
    year = request.args.get('year')
    season = request.args.get('season')
    query_string = f'select * from media where {f"season_year = {year}" if year else 'TRUE'} and {f"season = '{season}'" if season else 'TRUE'} and average_score is not null order by average_score desc limit 50 offset %s'
    db = create_db_connection()
    try:
        top_medias = db.fetchall(
            query_string,
            ((page - 1) * 50,))
        return jsonify({
            'status': 'success',
            'data': top_medias
        }), 200
    except Exception as e:
        return jsonify({
            'status': 'error',
            'message': str(e)
        }), 500
    finally:
        db.close()


if __name__ == '__main__':
    app.run(debug=True)
