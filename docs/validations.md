# Validations

## Listing validations

`conch validations` lists all available validations.  Each validation has an
UUID ID. In commands on a single validation, you may use the full idea or the
first 8 hexadecimal digits of the UUID (all characters before the first dash).
For example, `39cb3ab6-1963-4c9a-94ea-e2d9258d8be0` may be shortened to
`39cb3ab6`.

## Testing a validation with a device

`conch validation VALIDATION_ID test DEVICE_ID` tests a validation against a
device. The `DEVICE_ID` is a device's serial number. The command returns a
table of the validation results from running the validation. *This does not
store the validations results in the database*. This command is intended for
testing a new validation or new reporting agent code.

`conch validation VALIDATION_ID test DEVICE_ID` receives the input
data to validate from STDIN. The input data must be in JSON format. For
example, you can do any of the following to test a validation:

```bash
$ conch validation 39cb3ab6 test COFFEE < report_data.json

$ conch validation 39cb3ab6 test COFFEE <<EOF
{
	"power" : {
		"gigawatts" : 1.21
	}
}
EOF
```

## Creating and managing validation plans

Validation plans are collections of validations. Validations plans are is
executed during device report ingest.  Validations are independent and
un-ordered within a validation plan. A given validation may be in 0 or many
validation plans, and a validation plan may have 0, 1, or many validations
associated with it.

`conch validation-plans get` lists all available validation plans. Like
validation IDs, you may shorten the UUID to the first 8 characters in commands.

`conch validation-plan PLAN_ID validations` lists all validations associated
with the plan.

## Testing validation plans

`conch validation-plan VALIDATION_PLAN_ID test DEVICE_ID` tests a validation
plan against a device. Testing with a validation plan works identically to
testing with a single validation, except the input data is processed by all
validations in the validation plan. This command is useful for verifying a
device report can satisfy the schemas for all validations in a validation plans
when developing a reporting agent.

Like testing a validation, the command receives the JSON-formatted input data
from STDIN. Any of the following options work.

```bash
$ conch validation-plan 39cb3ab6 test COFFEE < report_data.json

$ conch validation-plan 39cb3ab6 test COFFEE <<EOF
{
	"power" : {
		"gigawatts" : 1.21
	}
}
EOF
```

