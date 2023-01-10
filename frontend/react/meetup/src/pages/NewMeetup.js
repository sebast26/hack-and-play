import NewMeetupForm from "../components/meetups/NewMeetupForm";

function NewMeetupPage() {
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
    )


  }

  return (
    <section>
      <h1>Add New Meetup</h1>
      <NewMeetupForm onAddMeetup={addMeetupHandler}/>
    </section>
  )
}

export default NewMeetupPage