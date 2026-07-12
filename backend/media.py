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


def generate_query(base_query, include_where=True, *filters):
    filters_string = ""
    params = []

    for clause, value in filters:
        if value is None:
            continue

        if filters_string != "":
            filters_string += "and "
        elif filters_string == "" and include_where:
            filters_string = "where "
        elif filters_string == "":
            filters_string = "and "

        filters_string += f"{clause} = %s "
        params.append(value)

    return base_query.replace("<<<where_clauses>>>", filters_string), params


@app.route('/')
def hello_world():  # put application's code here
    return 'Hello World!'


@app.route('/api/media/<int:id>')
def get_media(id):
    db = create_db_connection()
    try:
        media_data = db.fetchone('select *, to_json(titles) as titles from media where id = %s',(id,))
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
    media_type = request.args.get('media_type')

    db = create_db_connection()
    if page == 0: page = 1
    query_string = 'select *, to_json(titles) as titles from media <<<where_clauses>>> order by id limit 50 offset %s'
    query, params = generate_query(query_string, True, ("type", media_type))
    try:
        medias = db.fetchall(query, tuple(params) + ((page - 1) * 50,))
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

@app.route('/api/medias/batch', methods=['POST'])
def get_batch_medias():
    data = request.get_json()

    if not data or 'ids' not in data:
        return jsonify({
            'status': 'error',
            'message': 'No IDs provided'
        }), 400

    request_ids = data['ids']
    details_str = request.args.get('details', 'false').lower()

    if details_str not in ('true', 'false'):
        return jsonify({
            'status': 'error',
            'message': 'details must be true or false'
        }), 400

    if not isinstance(request_ids, list):
        return jsonify({
            'status': 'error',
            'message': 'ids must be a list'
        }), 400

    if not request_ids:
        return jsonify({
            'status': 'error',
            'message': 'ids cannot be empty',
        }), 400

    if not all(type(i) is int for i in request_ids):
        return jsonify({
            'status': 'error',
            'message': 'ids must be integers'
        }), 400

    if len(request_ids) > 30:
        return jsonify({
            'status': 'error',
            'message': 'a maximum of 30 ids is allowed',
        }), 400

    include_details = details_str == 'true'

    db = create_db_connection()
    try:
        query_string = f"select *, to_json(titles) as titles {',to_json(start_date) as start_date, to_json(end_date) as end_date, to_json(airing_schedule) as airing_schedule, to_json(recommendations) as recommendations, to_json(score_distribution) as score_distribution' if include_details else ''} from media {'left join media_details on media.id = media_details.id' if include_details else ''} where media.id = any(%s)"
        medias = db.fetchall(query_string, (request_ids,))
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
        media_data = db.fetchone('select *, to_json(titles) as titles, to_json(start_date) as start_date, to_json(end_date) as end_date, to_json(airing_schedule) as airing_schedule, to_json(recommendations) as recommendations, to_json(score_distribution) as score_distribution from media m left join media_details md on m.id = md.id where m.id = %s',
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
    media_type = request.args.get('media_type')

    query_string = 'select m.*, to_json(titles) as titles from media m left join media_details md on m.id = md.id <<<where_clauses>>> order by md.popularity desc limit 50 offset %s'
    query, params = generate_query(query_string, True, ("m.type", media_type), ("m.season_year", year), ("m.season", season))
    db = create_db_connection()
    try:
        popular_medias = db.fetchall(
            query,
            tuple(params) + ((page - 1) * 50,))
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
    media_type = request.args.get('media_type')

    query_string = 'select *, to_json(titles) as titles from media where average_score is not null <<<where_clauses>>> order by average_score desc limit 50 offset %s'
    query, params = generate_query(query_string, False, ("type", media_type), ("season_year", year), ("season", season))
    db = create_db_connection()
    try:
        top_medias = db.fetchall(
            query,
            tuple(params) + ((page - 1) * 50,))
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
