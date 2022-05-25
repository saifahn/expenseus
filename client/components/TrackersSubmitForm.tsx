export default function TrackersSubmitForm() {
  return (
    <div className="mt-6">
      <form className="border-4 p-6">
        <h3 className="text-lg font-semibold">Create Tracker</h3>
        <div className="mt-4">
          <label className="block font-semibold" htmlFor="name">
            Name
          </label>
          <input
            className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
            type="text"
            id="name"
          />
        </div>
        {/* TODO: make this a list of users that is populated by get all users */}
        <div className="mt-4">
          <label className="block font-semibold" htmlFor="participants">
            Participants
          </label>
          <input
            className="appearance-none w-full border rounded leading-tight focus:outline-none focus:ring py-2 px-3 mt-2"
            type="text"
            id="participants"
          />
        </div>
        <div className="mt-6 flex justify-end">
          <button className="border-4 py-2 px-4 rounded focus:outline-none focus:ring">
            Create tracker
          </button>
        </div>
      </form>
    </div>
  );
}
