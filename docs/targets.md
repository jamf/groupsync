# Adding new targets
Targets are services that can be synced *into*. To be able to do that, you need
to define how group members can be added/deleted to/from these services.

## Code things
1. First, [implement the Service interface for your thing and make sure it
   works](services.md).
2. Implement the remaining methods that consist the
   [Target interface](../services/target.go).
3. Add your target to the `TargetFromStr` function found in
   [target.go](../services/target.go).
4. In the `acquireIdentity` method of your new target, make sure there's
   logic for converting identities from services you might use as sources.