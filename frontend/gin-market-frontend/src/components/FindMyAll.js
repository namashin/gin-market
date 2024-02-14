import React, { useEffect, useState } from 'react';

const FindMyAll = () => {
    const [items, setItems] = useState([]);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchItems = async () => {
            try {
                const accessToken = localStorage.getItem('accessToken');
                const response = await fetch('http://localhost:8080/items/mine', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${accessToken}`
                    }
                });

                if (response.ok) {
                    const data = await response.json();
                    setItems(data.data);
                } else {
                    throw new Error('Failed to fetch items');
                }
            } catch (error) {
                console.error('Error fetching items:', error);
                setError('Failed to fetch items. Please try again later.');
            }
        };

        fetchItems();
    }, []);

    return (
        <div>
            <h2>My Items</h2>
            {error && <p>{error}</p>}
            <ul>
                <pre>{JSON.stringify(items, null, 2)}</pre>
            </ul>
        </div>
    );
};

export default FindMyAll;
