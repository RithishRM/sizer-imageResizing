import './App.css';
import { useState, useRef } from 'react';
import Title from './Components/Title';
import DropBox from './Components/DropBox';
import Resizer from './Components/Resizer';

function App() {
  const [details, setDetails] = useState(null);
  const [format, setFormat] = useState('png');  // Default format

  const handleFormatChange = (e) => {
    setFormat(e.target.value);
  };

  return (
    <div className="App">
      <Title/>
      <div className='flex justify-center'>
        <DropBox setting_img_details={setDetails}/>
        <Resizer get_img_details={details} set_img_details={setDetails} format={format}/>
        <div className="format-select">
          <label htmlFor="format">Choose format: </label>
          <select id="format" value={format} onChange={handleFormatChange}>
            <option value="png">PNG</option>
            <option value="jpeg">JPEG</option>
            <option value="gif">GIF</option>
          </select>
        </div>
      </div>
      {console.log(details, format)}
    </div>
  );
}

export default App;
