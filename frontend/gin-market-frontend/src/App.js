import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import FindAll from './components/FindAll';
import Signup from './components/Signup';
import Login from './components/Login';
import Create from './components/Create';
import Delete from './components/Delete';
import FindMyAll from "./components/FindMyAll";
import FindById from "./components/FindById";

const App = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<FindAll />} />
                <Route path="/signup" element={<Signup />} />
                <Route path="/login" element={<Login />} />
                <Route path="/create" element={<Create />} />
                <Route path="/delete/:id" element={<Delete />} />
                <Route path="/mine" element={<FindMyAll />} />
                <Route path="/find/:id" element={<FindById />} />
            </Routes>
        </BrowserRouter>
    );
};

export default App;
