import React, { useEffect, useState } from 'react';

const FindAll = () => {
    const [items, setItems] = useState([]);

    useEffect(() => {
        const backendURL = 'http://localhost:8080/items';

        fetch(backendURL)
            .then(response => response.json())
            .then(data => {
                // レスポンスデータが {"data": []*models.Item} 形式
                // レスポンスデータの配列を直接 items ステートにセット
                setItems(data.data);
            })
            .catch(error => {
                console.error('Error fetching data:', error);
            });
    }, []);

    return (
        <div>
            <h1>Item List</h1>
            <pre>{JSON.stringify(items, null, 2)}</pre>
        </div>
    );
};

export default FindAll;
