import NewMeetupForm from "../components/meetups/NewMeetupForm";
import {useNavigate} from 'react-router-dom'

function NewMeetupPage() {
  const navigate = useNavigate()

  function addMeetupHandler(meetupDate) {
    fetch(
      'http://example.com',
      {
        method: 'POST',
        body: JSON.stringify(meetupDate),
        headers: {
          'Content-Type': 'application/json'
        }
      }
    ).then(() => {
      navigate("/", { replace: true });
    }).catch(() => {
      navigate("/", { replace: true });
    })


  }

  return (
    <section>
      <h1>Add New Meetup</h1>
      <NewMeetupForm onAddMeetup={addMeetupHandler}/>
    </section>
  )
}

export default NewMeetupPage