select id, display_name, lifecycle_state
from oci.oci_core_internet_gateway
where display_name = '{{ output.display_name.value }}';