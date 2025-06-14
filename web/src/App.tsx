import './App.css';
import logo from './assets/stackjet logo.png';

function App() {
	return (
		<div
			style={{
				display: 'flex',
				justifyContent: 'center',
				alignItems: 'center',
				minHeight: '90vh',
			}}
		>
			<img src={logo} style={{ maxWidth: '700px', padding: '20px' }} />
		</div>
	);
}

export default App;
