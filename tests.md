## Workout Manager Tests
### Mocks in Workout Manager Tests

- **User and Peripheral Client Mocks**: These mocks simulate external services like user management and peripheral devices. They allow the tests to control the responses and behavior of these services, ensuring that the Workout Manager's interactions with these services can be tested independently of their actual implementations.

- **WorkoutStatsPublisher Mock**: This mock replaces the actual workout statistics publishing mechanism. It's used to verify if the Workout Manager is correctly publishing statistics, without needing to integrate with the real publishing system.

- **Postgres Repository**: Contrary to other components, the database interactions in the Workout Manager tests are not mocked. The tests interact with an actual Postgres repository, this was does for the ease of testing, the design allows us to plug a mock seamlessly.

### Tests - services_test.go
1. **TestWorkoutService_StartAndStop**: Validates the ability to start a new workout session and stop it correctly, ensuring that the peripheral device is bound and unbound properly.

2. **TestWorkoutService_StartWorkoutTwice**: Tests the scenario where a workout is started twice, expecting an error on the second attempt, and then stops the workout.

3. **TestWorkoutService_HardcoreModeNoShelter**: Ensures in hardcore mode, the shelter option is not available, and the workout can be started and stopped as expected.

4. **TestWorkoutService_WorkoutOptionsStartMultipleTimesStop**: Checks the behavior of starting and stopping workout options, including handling errors when stopping options that haven't been started.

5. **TestWorkoutService_UpdateDistanceTravelled**: Simulates the workout's distance tracking, updating the distance multiple times, and then verifies the total distance is correctly accumulated.

6. **TestWorkoutProcess_Shelters**: Validates that the shelter count in the workout data reflects the actual number of times shelter was taken.

7. **TestWorkoutService_HardcoreMode**: Tests if the shelter option is correctly omitted in hardcore mode and verifies the workout's start and stop functionality.

8. **TestWorkoutService_InitialWorkoutOptionsIfCardio**: Ensures that for a cardio profile, the workout options change based on the user's heart rate.
**It was not easily possible for us to mock the heartrate using Postman. The system does take into consideration the average heartrate for computing the workout options, however it may or may not be in the cardio zone to affect the options (heartrate > 70% of max heart rate). This case is tested in the integration tests in the Workout Manager.**

9. **TestWorkoutService_WorkoutOptionsIfCardio**: Tests that in cardio mode, the workout options adjust based on the user's actions during the workout.

10. **TestWorkoutService_InitialWorkoutOptionsIfStrength**: Checks the default order of workout options for a strength profile, ensuring it is 'fight' followed by 'escape'.

11. **TestWorkoutService_WorkoutOptionsIfStrength**: Verifies that in strength mode, the workout options adjust based on the user's actions during the workout.

12. **TestWorkoutService_DistanceToShelterUpdatesTest**: Tests the update mechanism for the distance to shelter in the workout options query, ensuring it changes as expected.

## Challenge Manager Tests
### Challenge Manager Service Tests - services_test.go

1. **TestChallengeService_CreateOrUpdateChallengeStats**: 
- **Mocking Approach**: The test employs mock data to simulate player workout statistics that are relevant to the challenge criteria. This includes mocking the distance covered, the number of enemies fought, and the number of times escaped, which align with the predefined challenge criteria like 'DistanceCovered' and 'FightMoreThanEscape'.
- **Test Description**: Tests the creation or updating of challenge stats and checks if the correct badges are awarded or not awarded based on different challenge criteria, like 'DistanceCovered' and 'FightMoreThanEscape'.

### Challenge Manager Domain Tests - domain_test.go
1. **TestChallenge_NewChallenge - check invalid criteria**: Verifies that creating a new challenge with invalid criteria results in an `ErrorInvalidCriteria`.

2. **TestChallenge_NewChallenge - Invalid challenge duration**: Checks that creating a challenge with an invalid duration (end time before start time) leads to an `ErrInvalidTime`.

## User Manager Tests
### User Manager Domain Tests - user_test.go

1. **TestPlayer_NewUser - Empty name validation**: Ensures that creating a user with an empty name results in an `ErrInvalidUserName` error.

2. **TestPlayer_NewUser - Empty email validation**: Verifies that creating a user with an empty email leads to an `ErrInvalidUserEmail` error.

3. **TestPlayer_NewUser - Empty dob validation**: Checks that creating a user with an empty date of birth (dob) causes an `ErrInvalidUserDOB` error.

4. **TestPlayer_NewUser - Valid user**: Confirms that creating a user with valid details does not result in any error.

### User Manager Domain Tests - player_test.go
1. **TestPlayer_NewPlayer - Empty weight validation**: Tests that creating a player with a weight of 0.0 triggers an `ErrInvalidPlayerWeight` error.

2. **TestPlayer_NewPlayer - Empty height validation**: Ensures that creating a player with a height of 0.0 leads to an `ErrInvalidPlayerHeight` error.

3. **TestPlayer_NewPlayer - Incorrect Preference validation**: Verifies that creating a player with an invalid preference (like 'jumping') results in an `ErrInvalidPlayerPreference` error.

4. **TestPlayer_NewPlayer - Empty zoneID validation**: Checks that creating a player with an empty (Nil) zoneID causes an `ErrInvalidZoneID` error.

## Peripheral Service Tests

### Mocks in Peripheral Service Tests

- **AMQP Publisher Mock**: Allows us to send live location data to Zone Manager and Workout Manager, without their actual instances.
- **Zone Service Client Mock**: Simulates the zone service client for testing interactions without a live zone service.

### Tests - services_test.go

1. **TestPeripheralService_Bind_Unbind**: Validates binding and unbinding a peripheral to a workout session, ensuring correct service operations.

2. **TestPeripheralService_UnbindWithoutBind**: Verifies handling of unbinding a peripheral that was never bound, expecting an error.

3. **TestDisconnectPeripheral_Exists**: Tests unbinding an existing peripheral, ensuring its removal from the repository.

4. **TestCreatePeripheral**: Confirms creating a new peripheral and its correct addition to the repository.

5. **TestCheckStatusByHRMId**: Ensures accurate reporting of a peripheral's status based on its HRM ID.

6. **TestBindPeripheral_NewPeripheral**: Tests creating new peripheral, verifying correct repository addition and status update.

7. **TestDisconnectPeripheral_NotExists**: Checks error handling when disconnect or unbind a non-existent peripheral.

8. **TestSetHeartRateReading**: Verifies updating and reflecting a bound peripheral's heart rate reading in the repository.

9. **TestGetHRMAvgReading**: Tests retrieving an average heart rate reading for a specific workout, validating data accuracy.

10. **TestGetHRMReading**: Confirms accurate and timely retrieval of current heart rate reading from a peripheral.

11. **TestGetHRMDevStatus**: Checks correct reporting of a peripheral's HRM device status.

12. **TestSetHRMDevStatusByHRMId**: Ensures accurate reflection of HRM status changes in the repository.

13. **TestSetGeoLocation**: Verifies setting and storing geolocation data for a peripheral accurately.

14. **TestGetGeoDevStatus**: Tests retrieval of a peripheral's geolocation device status, confirming accuracy.

15. **TestSetGeoDevStatus**: Checks the update and reflection of geolocation device status changes.

16. **TestGetGeoLocation**: Ensures accurate retrieval of a peripheral's current geolocation.

17. **TestGetLiveStatus**: Tests accurate retrieval of a peripheral's live status.

18. **TestSetLiveStatus**: Verifies correct recording of live status changes for a peripheral in the repository.

## Zone Manager Tests

### Mocks in Zone Manager Tests

- **Shelter Publisher Mock**: Simulates the AMQP publisher, allowing for testing of messaging functionalities without a real AMQP server - for sending out the shelter distances.
- **Postgres Repository**: Contrary to other components, the database interactions in the Workout Manager tests are not mocked. The tests interact with an actual Postgres repository, this was does for the ease of testing, the design allows us to plug a mock seamlessly.

### Tests - services_test.go

1. **TestZoneService_CreateTrail**: Confirms the creation of a trail, checking for a valid ID and error-free operation.

2. **TestZoneService_CreateShelter**: Verifies the creation of a shelter, ensuring it is properly registered with a unique ID.

3. **TestZoneService_UpdateTrail**: Tests updating an existing trail, ensuring the changes are accurately reflected.

4. **TestZoneService_DeleteTrail**: Confirms the ability to delete a trail, verifying its removal from the system.

5. **TestZoneService_GetTrailByID**: Checks the retrieval of a trail by its ID, confirming accurate data fetching.
