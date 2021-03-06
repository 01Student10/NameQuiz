import sirv from 'sirv';
import polka from 'polka';
import compression from 'compression';
import * as sapper from '@sapper/server';
import session from 'express-session';
import bodyParser from 'body-parser';
import sessionFileStore from 'session-file-store';


const FileStore = sessionFileStore(session)


const { PORT, NODE_ENV } = process.env;
const dev = NODE_ENV === 'development';

polka()// You can also use Express
	.use(bodyParser.json())
	.use(session({
		secret: 'secret',
		resave: false,
		saveUninitialized: true,
		cookie: {
			maxAge: 999999999
		},
		store: new FileStore({
			path: process.env.NOW ? `/tmp/sessions` : `.sessions`
		})
	})
	)
	.use(
		compression({ threshold: 0 }),
		sirv('static', { dev }),
		sapper.middleware({
			session: req => ({
				user: req.session && req.session.user
			})
		})
	)
	.listen(PORT, err => {
		if (err) console.log('error', err);
	});
