import { useState } from 'react';
import { gql } from '@apollo/client';
import { useQuery, useMutation } from '@apollo/client/react';
import './App.css';

const GET_DROP = gql`
  query GetDrop($id: ID!) {
    getDrop(id: $id) {
      id
      title
      totalAvailable
      claimed
    }
  }
`;

const CLAIM_STICKER = gql`
  mutation ClaimSticker($dropId: ID!, $email: String!) {
    claimSticker(dropId: $dropId, email: $email)
  }
`;

interface DropData {
  getDrop: {
    id: string;
    title: string;
    totalAvailable: number;
    claimed: number;
  };
}

function App() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');

  const generateRandomEmail = () => {
    const randomString = Math.random().toString(36).substring(2, 8);
    setEmail(`tester_${randomString}@mule.com`);
    setMessage('');
  };

  const { loading, error, data, refetch } = useQuery<DropData>(GET_DROP, {
    variables: { id: "1" }
  });

  const [claimSticker, { loading: claimLoading }] = useMutation(CLAIM_STICKER, {
    onError: (err: any) => {
      setMessage(`${err.message}`);
    },
    onCompleted: () => {
      setMessage("Congrats! - Check your inbox!");
      refetch(); // reloads the drop data so Counter can jump +1!
    }
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault(); // prevent site relaod
    setMessage('');
    
    if (!email) {
      setMessage("Please enter a E-Mail address.");
      return;
    }

    claimSticker({ variables: { dropId: "1", email: email } });
  };

  if (loading) return <p>Loading Drop...</p>;
  if (error) return <p>Error: {error.message}</p>;

  if (!data) return <p>No Drop Data found!</p>;

  const drop = data.getDrop;
  const remaining = drop.totalAvailable - drop.claimed;

  return (
    <div className="card">
      <h1>Sticker Drop</h1>
      <h2>{drop.title}</h2>
      
      <div className="progress-container">
        <p><strong>{remaining}</strong> of {drop.totalAvailable} left!</p>
        <progress value={drop.claimed} max={drop.totalAvailable}></progress>
      </div>

      <form onSubmit={handleSubmit} className="claim-form">
        <div className="input-group">
          <input 
            type="email" 
            placeholder="your@email.com" 
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={claimLoading || remaining === 0}
          />
          <button 
            type="button" 
            className="random-btn" 
            onClick={generateRandomEmail}
            disabled={claimLoading || remaining === 0}
            title="generate random E-Mail"
          >
          Generate
          </button>
        </div>
        <button type="submit" disabled={claimLoading || remaining === 0}>
          {claimLoading ? 'Validating...' : 'Securing Sticker!'}
        </button>
      </form>

      <p className="message">{message}</p>
    </div>
  );
}

export default App;